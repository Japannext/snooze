#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Module for managing the file based database backend (with TinyDB)'''

import re
from copy import deepcopy
from datetime import datetime
from functools import reduce
from logging import getLogger
from pathlib import Path
from threading import Lock
from typing import List, Tuple

import uuid
from bson.json_util import dumps
from tinydb import TinyDB, Query as BaseQuery, where
from tinydb.operations import add

from snooze.db.database import Database, wrap_exception
from snooze.utils.config import FileConfig
from snooze.utils.functions import dig, to_tuple, flatten

log = getLogger('snooze.db.file')
DEFAULT_PAGINATION = {
    'orderby': '',
    'page_number': 1,
    'nb_per_page': 0,
    'asc': True,
    'only_one': False,
}

class OperationNotSupported(Exception):
    '''Raised when the search operator is not supported'''

mutex = Lock()

class Query(BaseQuery):
    def test_root(self, func, *args):
        return self._generate_test(
            lambda value: func(value, *args),
            ('test', self._path, func, args),
            allow_empty_path=True
        )


def test_contains(array, value):
    if not isinstance(array, list):
        array = [array]
    for val in value:
        if isinstance(val, str):
            reg = re.compile(val, flags=re.IGNORECASE)
            for record in array:
                if reg.search(record):
                    return True
        else:
            for record in array:
                if str(val) in str(record):
                    return True
    return False

def test_search(dic, value):
    return value in str(dic)

def append_many(fields):
    def transform(doc):
        for field, val in fields.items():
            doc[field] += val
    return transform

def prepend_many(fields):
    def transform(doc):
        for field, val in fields.items():
            doc[field] = val + doc[field]
    return transform

def set_many(fields):
    def transform(doc):
        for field, val in fields.items():
            doc[field] = val
    return transform

def remove_many(fields):
    def transform(doc):
        for field, val in fields.items():
            doc[field] = [i for i in doc[field] if i not in val]
    return transform

def overwrite(data, field):
    def transform(doc):
        if doc[field] in data:
            for k,v in data[doc[field]].items():
                doc[k] = v
    return transform

class BackendDB(Database):
    '''Backend database based on local file (TinyDB)'''

    name = 'file'

    def __init__(self, config: FileConfig):
        self.db = TinyDB(config.path)
        log.debug("Initialized TinyDB at path %s", config.path)
        log.debug("db: %s", self.db)
        log.debug("List of collections: %s", self.db.tables())

    def create_index(self, collection, fields):
        pass

    def cleanup_timeout(self, collection):
        mutex.acquire()
        now = datetime.now().timestamp()
        aggregate_results = self.db.table(collection).search(Query().ttl >= 0)
        aggregate_results = [{'_id': doc.doc_id, 'timeout': doc['ttl'] + doc['date_epoch']} for doc in aggregate_results]
        aggregate_results = [doc for doc in aggregate_results if doc['timeout'] <= now]
        res = self.delete_aggregates(collection, aggregate_results)
        mutex.release()
        return res

    def cleanup_comments(self):
        '''Delete comments which record doesn't exist anymore'''
        with mutex:
            record_uids = [doc['uid'] for doc in self.db.table('record').all()]
            aggregate_results = self.db.table('comment').search(~ (Query()['record_uid'].one_of(record_uids)))
            aggregate_results = [{'_id': doc.doc_id} for doc in aggregate_results]
            res = self.delete_aggregates('comment', aggregate_results)
            return res

    def cleanup_orphans(self, collection: str) -> int:
        '''Delete objects which one of their ancestors does not exist anymore'''
        with mutex:
            parents = set(flatten([doc.get('parents', []) for doc in self.db.table(collection).all()]))
        if len(parents) == 0:
           return 0
        to_delete = []
        for parent in parents:
            if not self.get_one(collection, {'uid': parent}):
                to_delete.append(parent)
        total = 0
        with mutex:
            total = len(self.db.table(collection).remove(Query().parents.any(to_delete)))
        return total

    def cleanup_audit_logs(self, interval):
        mutex.acquire()
        audits = self.db.table('audit').all()
        audit_dict = {}
        now = datetime.now().timestamp()
        threshold_date = now - interval
        for audit in audits:
            oid = audit.get('object_id')
            timestamp = audit_dict.get(oid, {}).get('date_epoch', 0)

            # Remember latest audit log by object_id
            if oid and audit.get('date_epoch', 0) >= timestamp:
                audit_dict[oid] = audit

        oids = [
            audit['object_id']
            for object_id, audit in audit_dict.items()
            if audit.get('action') == 'deleted' \
            and audit.get('date_epoch') < threshold_date
        ]
        log.debug("Found audit log to remove for %d objects", len(oids))
        doc_ids = [
            obj.doc_id
            for oid in oids
            for obj in self.db.table('audit').search(Query().object_id == oid)
        ]
        log.debug("Found %d audit logs to remove", len(doc_ids))
        self.db.table('audit').remove(doc_ids=doc_ids)

        mutex.release()

    def renumber_field(self, collection, field):
        '''Renumber field by ascending order'''
        log.info("Reordering field '%s' in collection %s", field, collection)
        data = self.search(collection, pagination={'orderby': field, 'asc': True})['data']
        rows_by_uid = {}
        for idx, row in enumerate(data):
            row[field] = idx
            rows_by_uid[row['uid']] = row
        self.db.table(collection).update(overwrite(rows_by_uid, 'uid'))
        log.info("Field '%s' renumbering on collection %s: Success", field, collection)

    def delete_aggregates(self, collection, aggregate_results):
        ids = [doc['_id'] for doc in aggregate_results]
        deleted_count = 0
        if ids:
            deleted_results = self.db.table(collection).remove(doc_ids=ids)
            deleted_count = len(deleted_results)
            log.debug('Removed %d documents in %s', deleted_count, collection)
        return deleted_count

    @wrap_exception
    def write(self, collection, obj, primary=None, duplicate_policy='update', update_time=True, constant=None):
        mutex.acquire()
        added = []
        updated = []
        replaced = []
        rejected = []
        obj_copy = []
        tobj = obj
        add_obj = False
        table = self.db.table(collection)
        tobj = deepcopy(obj)
        if not isinstance(tobj, list):
            tobj = [tobj]
        if primary:
            if isinstance(primary , str):
                primary = primary.split(',')
        if constant:
            if isinstance(constant , str):
                constant = constant.split(',')
        for o in tobj:
            primary_docs = None
            old = {}
            o.pop('_old', None)
            if update_time:
                o['date_epoch'] = datetime.now().timestamp()
            if primary and all(dig(o, *p.split('.')) for p in primary):
                primary_query = map(lambda a: dig(Query(), *a.split('.')) == dig(o, *a.split('.')), primary)
                primary_query = reduce(lambda a, b: a & b, primary_query)
                primary_docs = table.search(primary_query)
                if primary_docs:
                    log.debug('Documents with same primary %s: %s', primary, primary_docs[0].doc_id)
            if 'uid' in o:
                query = Query()
                docs = table.search(query.uid == o['uid'])
                if docs:
                    doc = docs[0]
                    doc_id = doc.doc_id
                    old = doc
                    log.debug('Found: %s', doc_id)
                    if primary_docs and doc_id != primary_docs[0].doc_id:
                        error_message = f"Found another document with same primary {primary}: {primary_docs}. " \
                            "Since UID is different, cannot update"
                        log.error(error_message)
                        o['error'] = error_message
                        rejected.append(o)
                    elif constant and any(doc.get(c, '') != o.get(c) for c in constant):
                        error_message = f"Found a document with existing uid {o['uid']} but different constant " \
                            f"values: {constant}. Since UID is different, cannot update"
                        log.error(error_message)
                        o['error'] = error_message
                        rejected.append(o)
                    elif duplicate_policy == 'replace':
                        log.debug('Replacing with: %s', doc_id)
                        table.remove(doc_ids=[doc_id])
                        table.insert(o)
                        replaced.append(o)
                    else:
                        log.debug('Updating with: %s', doc_id)
                        table.update(o, doc_ids=[doc_id])
                        updated.append(o)
                else:
                    error_message = f"UID {o['uid']} not found. Skipping..."
                    log.error(error_message)
                    o['error'] = error_message
                    rejected.append(o)
            elif primary:
                if primary_docs:
                    doc = primary_docs[0]
                    doc_id = doc.doc_id
                    old = doc
                    if constant and any(doc.get(c, '') != o.get(c) for c in constant):
                        error_message = f"Found a document with existing primary {primary} but different "\
                            f"constant values: {constant}. Since UID is different, cannot update"
                        log.error(error_message)
                        o['error'] = error_message
                        rejected.append(o)
                    else:
                        log.debug('Evaluating duplicate policy: %s', duplicate_policy)
                        if duplicate_policy == 'insert':
                            add_obj = True
                        elif duplicate_policy == 'reject':
                            error_message = "Another object exist with the same {primary}"
                            o['error'] = error_message
                            rejected.append(o)
                        elif duplicate_policy == 'replace':
                            log.debug('Replace with: %s', doc_id)
                            table.remove(doc_ids=[doc_id])
                            if 'uid' in doc:
                                o['uid'] = doc['uid']
                            table.insert(o)
                            replaced.append(o)
                        else:
                            log.debug('Update with: %s', doc_id)
                            table.update(o, doc_ids=[doc_id])
                            updated.append(o)
                else:
                    log.debug("Could not find document with primary {}. Inserting instead".format(primary))
                    add_obj = True
            else:
                add_obj = True
            if add_obj:
                obj_copy.append(o)
                obj_copy[-1]['uid'] = str(uuid.uuid4())
                added.append(o)
                add_obj = False
                log.debug("In %s, inserting %s", collection, o.get('uid', ''))
            if old:
                o['_old'] = old
        if len(obj_copy) > 0:
            table.insert_multiple(obj_copy)
        mutex.release()
        return {
            'data': {
                'added': deepcopy(added),
                'updated': deepcopy(updated),
                'replaced': deepcopy(replaced),
                'rejected': deepcopy(rejected),
            },
        }

    @wrap_exception
    def get_one(self, collection: str, search: dict):
        '''Return the first element found based on a search dict'''
        with mutex:
            queries = []
            for key, value in search.items():
                query = dig(Query(), *key.split('.')) == value
                queries.append(query)
            search_query = reduce(lambda a, b: a & b, queries)
            results = self.db.table(collection).search(search_query)
            if results:
                return results[0]
            else:
                return None

    @wrap_exception
    def replace_one(self, collection: str, search: dict, obj: dict, update_time: bool = True):
        with mutex:
            query = Query()
            new_obj = dict(obj)
            queries = []
            for key, value in search.items():
                new_obj[key] = value
                query = dig(Query(), *key.split('.')) == value
                queries.append(query)
            if update_time:
                new_obj['date_epoch'] = datetime.now().timestamp()
            search_query = reduce(lambda a, b: a & b, queries)
            search = self.db.table(collection).search(search_query)
            if search:
                self.db.table(collection).remove(doc_ids=[search[0].doc_id])
            self.db.table(collection).insert(new_obj)
            return len(list(search))

    @wrap_exception
    def update_one(self, collection: str, uid: str, obj: dict, update_time: bool = True):
        with mutex:
            new_obj = dict(obj)
            if update_time:
                new_obj['date_epoch'] = datetime.now().timestamp()
            self.db.table(collection).upsert(new_obj, where('uid') == uid)

    @wrap_exception
    def bulk_increment(self, collection: str, updates: List[Tuple[dict, dict]], upsert: bool = False):
        '''Perform a bulk update of increments. Each update should be a tuple of search and update'''
        with mutex:
            for search, update in updates:
                queries = []
                for key, value in search.items():
                    if isinstance(value, datetime):
                        value = int((value.timestamp() // 3600) * 3600)
                        search[key] = value
                    query = dig(Query(), *key.split('.')) == value
                    queries.append(query)
                search_query = reduce(lambda a, b: a & b, queries)
                update.pop('_id', None)
                for key, value in update.items():
                    self.db.table(collection).update(add(key, value), search_query)
                if upsert and not self.db.table(collection).search(search_query):
                    self.db.table(collection).insert({**search, **update})

    @wrap_exception
    def inc(self, collection, field, labels={}):
        now = int((datetime.now().timestamp() // 3600) * 3600)
        mutex.acquire()
        table = self.db.table(collection)
        query = Query()
        keys = []
        added = []
        updated = []
        if labels:
            for key, value in labels.items():
                keys.append(f"{field}__{key}__{value}")
        else:
            keys.append(field)
        for key in keys:
            result = table.search((query.date == now) & (query.key == key))
            if result:
                result = result[0]
                result['value'] = result.get('value', 0) + 1
                table.update(result, doc_ids=[result.doc_id])
                updated.append(deepcopy(result))
            else:
                result = {'date': now, 'type': 'counter', 'key': key}
                result['value'] = 1
                table.insert(result)
                added.append(deepcopy(result))
        mutex.release()
        return {'data': {'added': added, 'updated': updated}}

    def inc_many(self, collection: str, field: str, condition = None, value: int = 1):
        if condition is None:
            condition = []
        tinydb_search = self.convert(condition)
        total = 0
        if collection in self.db.tables():
            table = self.db.table(collection)
            total = table.count(tinydb_search)
            table.update(add(field, value), tinydb_search)
        return total

    def update_with_operation(self, collection, operation, condition=[]):
        tinydb_search = self.convert(condition)
        total = 0
        mutex.acquire()
        if collection in self.db.tables():
            try:
                total = len(self.db.table(collection).update(operation, tinydb_search))
            except Exception as err:
                log.exception(err)
        mutex.release()
        log.debug("Updated %d document(s)", total)
        return total

    def set_fields(self, collection, fields, condition=[]):
        log.debug("Update collection '%s' with fields '%s' based on the following search '%s'", collection, fields, condition)
        return self.update_with_operation(collection, set_many(fields), condition)

    def append_list(self, collection, fields, condition=[]):
        log.debug("Append to collection '%s' fields '%s' based on the following search: '%s'", collection, fields, condition)
        return self.update_with_operation(collection, append_many(fields), condition)

    def prepend_list(self, collection, fields, condition=[]):
        log.debug("Prepend to collection '%s' fields '%s' based on the following search: '%s'", collection, fields, condition)
        return self.update_with_operation(collection, prepend_many(fields), condition)

    def remove_list(self, collection, fields, condition=[]):
        log.debug("Remove from collection '%s' fields '%s' based on the following search: '%s'", collection, fields, condition)
        return self.update_with_operation(collection, remove_many(fields), condition)

    def compute_stats(self, collection, date_from, date_until, groupby='hour'):
        log.debug("Compute metrics on `%s` from %s until %s grouped by %s", collection, date_from, date_until, groupby)
        date_from = date_from.replace(minute=0, second=0, microsecond=0)
        mutex.acquire()
        if collection not in self.db.tables():
            log.debug("Compute stats: collection %s does not exist", collection)
            mutex.release()
            return {'data': [], 'count': 0}
        if groupby == 'hour':
            date_format = '%Y-%m-%dT%H:00%z'
        elif groupby == 'day':
            date_format = '%Y-%m-%dT00:00%z'
        elif groupby == 'month':
            date_format = '%Y-%m-01T00:00%z'
        elif groupby == 'year':
            date_format = '%Y-01-01T00:00%z'
        elif groupby == 'week':
            date_format = '%Y-%VT00:00%z'
        elif groupby == 'weekday':
            date_format = '%u'
        else:
            date_format = '%Y-%m-%dT%H:00%z'
        date_from = date_from.timestamp()
        date_until = date_until.timestamp()
        table = self.db.table(collection)
        results = table.search((Query().date >= date_from) & (Query().date <= date_until))
        if len(results) == 0:
            log.debug("Compute stats: No data found within time interval")
            mutex.release()
            return {'data': [], 'count': 0}
        groups = {}
        res = []
        for doc in results:
            date_range = datetime.fromtimestamp(doc['date']).astimezone().strftime(date_format)
            if date_range not in groups:
                groups[date_range] = {doc['key']: {'value': 0}}
            elif doc['key'] not in groups[date_range]:
                groups[date_range][doc['key']] = {'value': 0}
            groups[date_range][doc['key']]['value'] += doc['value']
        for date, value in groups.items():
            entry = {'_id': date, 'data': []}
            for key, doc in value.items():
                entry['data'].append({'key': key, 'value': doc['value']})
            res.append(entry)
        results_agg = sorted(res, key=lambda d: d['_id'])
        count = len(results_agg)
        log.debug("Compute stats: Got %s results", count)
        mutex.release()
        return {'data': results_agg, 'count': count}

    @wrap_exception
    def search(self, collection, condition=[], **pagination):
        mutex.acquire()
        pagination = {**DEFAULT_PAGINATION, **pagination}
        orderby = pagination['orderby']
        page_number = pagination['page_number']
        nb_per_page = pagination['nb_per_page']
        asc = pagination['asc']
        only_one = pagination['only_one']
        tinydb_search = self.convert(condition)
        if collection in self.db.tables():
            table = self.db.table(collection)
            if tinydb_search:
                results = table.search(tinydb_search)
            else:
                results = table.all()
            if len(orderby) > 0 and all(dig(res, *orderby.split('.')) for res in list(results)):
                results = sorted(list(results),  key=lambda x: reduce(lambda c, k: c.get(k, {}), orderby.split('.'), x))
            if not asc:
                results = list(reversed(results))
            if only_one:
                total = 1
                results = results[:1]
            else:
                total = len(results)
                if nb_per_page > 0:
                    from_el = max((page_number-1)*nb_per_page, 0)
                    to_el = page_number*nb_per_page
                else:
                    from_el = None
                    to_el = None
                results = results[from_el:to_el]
            log.debug("Found %d result(s) for search %s in collection %s. Pagination: %s",
                total, tinydb_search, collection, pagination)
            mutex.release()
            return {'data': deepcopy(results), 'count': total}
        else:
            log.warning("Cannot find collection %s", collection)
            mutex.release()
            return {'data': [], 'count': 0}

    @wrap_exception
    def delete(self, collection, condition=[], force=False):
        mutex.acquire()
        tinydb_search = self.convert(condition)
        if collection in self.db.tables():
            table = self.db.table(collection)
            if len(condition) == 0 and not force:
                results_count = 0
                log.debug("Too dangerous to delete everything. Aborting")
            else:
                if len(condition) == 0:
                    results_count = len(table)
                    results = table.truncate()
                else:
                    results = table.remove(tinydb_search)
                    results_count = len(results)
                log.debug("Found %d item(s) to delete in collection %s for search %s",
                    results_count, collection, tinydb_search)
            mutex.release()
            return {'data': [], 'count': results_count}
        else:
            mutex.release()
            log.error("Cannot find collection %s", collection)
            return {'data': 0}

    def drop(self, collection):
        if collection in self.db.tables():
            self.db.drop_table(collection)

    def convert(self, array):
        """
        Convert `Condition` type from snooze.utils
        to Mongodb compatible type of search
        """
        if not array:
            return None
        operation, *args = array
        if operation == 'AND':
            arguments = list(map(self.convert, args))
            return_obj = reduce(lambda a, b: a & b, arguments)
        elif operation == 'OR':
            arguments = list(map(self.convert, args))
            return_obj = reduce(lambda a, b: a | b, arguments)
        elif operation == 'NOT':
            arg = self.convert(args[0])
            return_obj = ~ arg
        elif operation == '=':
            key, value = args
            return_obj = dig(Query(), *key.split('.')) == value
        elif operation == '!=':
            key, value = args
            return_obj = dig(Query(), *key.split('.')) != value
        elif operation == '>':
            key, value = args
            return_obj = dig(Query(), *key.split('.')) > value
        elif operation == '>=':
            key, value = args
            return_obj = dig(Query(), *key.split('.')) >= value
        elif operation == '<':
            key, value = args
            return_obj = dig(Query(), *key.split('.')) < value
        elif operation == '<=':
            key, value = args
            return_obj = dig(Query(), *key.split('.')) <= value
        elif operation == 'MATCHES':
            key, value = args
            return_obj = dig(Query(), *key.split('.')).search(value, flags=re.IGNORECASE)
        elif operation == 'EXISTS':
            return_obj = dig(Query(), *args[0].split('.')).exists()
        elif operation == 'CONTAINS':
            key, value = args
            if not isinstance(value, list):
                value = [value]
            for val in value:
                return_obj = dig(Query(), *key.split('.')).test(test_contains, to_tuple(value))
        elif operation == 'IN':
            key, value = args
            converted = False
            if not isinstance(key, list):
                key = [key]
            else:
                try:
                    saved_key = key
                    key = self.convert(key)
                    converted = True
                except Exception:
                    key = saved_key
            if converted:
                return_obj = dig(Query(), *value.split('.')).any(key)
            else:
                test_eq = lambda s, v: s == v
                eq_list = list(map(lambda a: dig(Query(), *value.split('.')).test(test_eq, a), key))
                return_obj = reduce(lambda a, b: a | b, eq_list) | dig(Query(), *value.split('.')).any(key)
        elif operation == 'SEARCH':
            arg = args[0]
            try:
                return_obj = Query().test_root(test_search, to_tuple(arg))
            except Exception as err:
                log.exception(err)
                raise OperationNotSupported(operation)
        else:
            raise OperationNotSupported(operation)
        return return_obj

    def backup(self, backup_path, backup_exclude=[]):
        collections = [c for c in self.db.tables() if c not in backup_exclude]
        log.debug('Starting backup of %s', collections)
        succeeded = []
        for i, collection_name in enumerate(collections):
            try:
                collection = self.db.table(collections[i]).all()
                jsonpath = Path(backup_path) / (collection_name + '.json')
                with jsonpath.open("wb") as jsonfile:
                    jsonfile.write(dumps(collection).encode())
                    succeeded.append(collection_name)
            except Exception as err:
                log.error('Backup of %s failed', collection_name)
                log.exception(err)
        log.info('Backup of %s succeeded', succeeded)
