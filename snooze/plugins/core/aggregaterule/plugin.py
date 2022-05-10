#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Plugin for aggregating records by a group of keys, thus avoiding users
from getting spammed for the same alert'''

import hashlib
from logging import getLogger
from datetime import datetime, timedelta
from typing import List

from pydantic import Field

from snooze.database import Pagination
from snooze.plugins.core import ApiPlugin, BaseProcessPlugin
from snooze.plugins.core.basic.plugin import Plugin, AbortAndUpdate
from snooze.utils.functions import dig
from snooze.utils.condition import Condition, AlwaysTrue
from snooze.utils.model import Record, ApiModel, MongodbMetadata, Comment

log = getLogger('snooze.aggregaterule')

class AggregateRule(ApiModel):
    _mongodb = MongodbMetadata(collection='aggregaterule')

    enabled: bool = True
    name: str = ''
    fields : List[str] = Field(default_factory=list)
    watch : List[str] = Field(default_factory=list)
    throttle: timedelta = timedelta(minutes=10)
    flapping: timedelta = timedelta(minutes=3)
    comment: str = ''
    condition: Condition = Field(default_factory=AlwaysTrue, discriminator='type')

    def match(self, record: Record) -> bool:
        '''Check if a record matched this aggregate's rule condition'''
        match = self.condition.match(record)
        if match:
            record.core.aggregate = self.name
        return match

    def hash(self, record: Record) -> str:
        '''Return the hash computed by an aggregaterule'''
        record_fields_strings = [f"{field}={dig(record, *field.split('.'))}" for field in self.fields]
        string_to_hash = self.name + '.'.join(record_fields_strings)
        return hashlib.md5(string_to_hash).hexdigest()

def default_hash(record: Record) -> str:
    '''Return a hash computed from any relevant field'''
    if not 'raw' in record:
        return hashlib.md5(repr(sorted(record.items())).encode('utf-8')).hexdigest()
    if isinstance(record['raw'], dict):
        return hashlib.md5(repr(sorted(record['raw'].items())).encode('utf-8')).hexdigest()
    elif isinstance(record['raw'], list):
        return hashlib.md5(repr(sorted(record['raw'])).encode('utf-8')).hexdigest()
    return hashlib.md5(record['raw'].encode('utf-8')).hexdigest()

def update_record(aggregate: Record, record: Record) -> Record:
    '''Update a record with the aggregate information'''
    #record = dict(list(aggregate.items()) + list(record.items()))
    record.uid = aggregate.uid
    record.state = aggregate.state
    record.core.duplicates = aggregate.core.duplicates + 1
    record.core.date_epoch = aggregate.core.date_epoch
    record.core.notification_from = None
    if aggregate.ttl < 0:
        record.ttl = aggregate.ttl

    return record

class AggregaterulePlugin(ApiPlugin, BaseProcessPlugin):
    '''The aggregate rule plugin'''
    aggregate_rules: List[AggregateRule]

    def __init__(self, core: 'Core'):
        ApiPlugin.__init__(self, core)
        pagination = Pagination()
        self.aggregate_rules = self.database[AggregateRule].search(pagination=pagination)

    def reload(self):
        '''Reload the plugin (data, setup, etc)'''
        pagination = Pagination()
        self.aggregate_rules = self.database[AggregateRule].search(pagination=pagination)

    def process(self, record: Record) -> Record:
        """Process the record against a list of aggregate rules"""
        log.debug("Processing record against aggregate rules")
        for aggrule in self.aggregate_rules:
            if aggrule.enabled and aggrule.match(record):
                record.hash = aggrule.hash(record)
                log.debug("Aggregate rule %s matched record: %s", self.name, record.hash)
                record = self.match_aggregate(record, aggrule.throttle, aggrule.flapping, aggrule.watch, aggrule.name)
                break
        else:
            log.debug("Record %s could not match any aggregate rule, assigning a default aggregate", record)
            record.hash = default_hash(record)
            record.aggregate = 'default'
            record = self.match_aggregate(record)

        return record

    def comment(self, aggrule: AggregateRule, aggregate: Record, record: Record) -> Comment:
        '''Return the comment associated with the aggregate update'''
        comment = Comment(record_uid=aggregate.uid, auto=True)
        if record_state == 'close':
            self.core.stats.inc('alert_closed', {'name': aggrule_name})
            if record.get('state') != 'close':
                log.debug("OK received, closing alert")
                aggregate_severity = aggregate.get('severity', 'unknown')
                record_severity = record.get('severity', 'unknown')
                comment.message = f"Auto closed: Severity {aggregate_severity} => {record_severity}"
                comment.type = 'close'
                record.state = 'close'
                self.database[Comment].create(comment)
                record.core.comment_count = aggregate.comment_count + 1
                return record
            else:
                log.debug("OK received but the alert is already closed, discarding")
                raise AbortAndUpdate(record)
        watched_fields = []
        for watched_field in watch:
            aggregate_field = dig(aggregate, *watched_field.split('.'))
            record_field = dig(record, *watched_field.split('.'))
            log.debug("Watched field %s: compare %s and %s", watched_field, record_field, aggregate_field)
            if record_field != aggregate_field:
                watched_fields.append({'name': watched_field, 'old': aggregate_field, 'new': record_field})
        if watched_fields:
            log.debug("Alert %s Found updated fields from watchlist: %s", str(record['hash']), watched_fields)
            append_txt = []
            for watch_field in watched_fields:
                append_txt.append("{watch_field['name']} ({watch_field['old']} => {watch_field['new']})")
            if record.get('state') == 'close':
                comment['message'] = 'Auto re-opened from watchlist: {}'.format(', '.join(append_txt))
                comment['type'] = 'open'
                record['state'] = 'open'
            elif record.get('state') == 'ack':
                comment['message'] = 'Auto re-escalated from watchlist: {}'.format(', '.join(append_txt))
                comment['type'] = 'esc'
                record['state'] = 'esc'
            else:
                comment['message'] = 'New escalation from watchlist: {}'.format(', '.join(append_txt))
                comment['type'] = 'comment'
            self.db.write('comment', comment)
            record['flapping_countdown'] = aggregate.get('flapping_countdown', flapping) - 1
        elif record.get('state') == 'close':
            comment['message'] = 'Auto re-opened'
            comment['type'] = 'open'
            record['state'] = 'open'
            self.db.write('comment', comment)
            record['flapping_countdown'] = aggregate.get('flapping_countdown', flapping) - 1
        elif (throttle < 0) or (now.timestamp() - aggregate.get('date_epoch', 0) < throttle):
            log.debug("Alert %s Time within throttle %s range, discarding", record['hash'], throttle)
            self.core.stats.inc('alert_throttled', {'name': aggrule_name})
            raise AbortAndUpdate(record)
        else:
            if record.get('state') == 'ack':
                comment['type'] = 'esc'
                record['state'] = 'esc'
            else:
                comment['type'] = 'comment'
            comment['message'] = 'New escalation'
            self.db.write('comment', comment)
            record.pop('flapping_countdown', '')
        record['comment_count'] = aggregate.get('comment_count', 0) + 1

    def throttle(self, aggrule: AggregateRule, record: Record):
        '''Check if a given aggregate rule shold throttle the record, and modify the record
        if needed.'''
        log.debug("Checking if an aggregate with hash %s can be found", record.hash)
        #aggregate = self.endpoint.get(hash=record.hash)
        aggregate = self.endpoint.get(Equals(field='hash', value=record.hash))
        if not aggregate:
            log.debug("Not found, creating a new aggregate")
            record.core.duplicates = 1
            record.core.snoozed = None
            record.core.notifications = []
            return record
        log.debug("Found record hash %s, updating it with the record infos", record['hash'])
        now = datetime.now()
        record = update_record(aggregate, record)
        comment = self.comment(aggrule, aggregate, record)
        if record.core.flapping_countdown < 0:
            log.debug("Alert %s is flapping, discarding", record.hash)
            raise AbortAndUpdate(record)
        else:
            log.debug("Not found, creating a new aggregate")
            record['duplicates'] = 1
        record.pop('snoozed', '')
        record.pop('notifications', '')
        return record

    def match_aggregate(self, record, throttle=10, flapping=2, watch=[], aggrule_name='default'):
        '''Attempt to match an aggregate with a record, and throttle the record if it does'''
        log.debug("Checking if an aggregate with hash %s can be found", record['hash'])
        aggregate_result = self.db.search('record', ['=', 'hash', record['hash']])
        aggregate = self.endpoint.get(Equals(field='hash', value=record.hash))
        if aggregate:
            log.debug("Found record hash %s, updating it with the record infos", record['hash'])
            now = datetime.now()
            record = dict(list(aggregate.items()) + list(record.items()))
            record_state = record.get('state', '')
            record.uid = aggregate.uid
            record.state = aggregate.state
            record['duplicates'] = aggregate.get('duplicates', 0) + 1
            record['date_epoch'] = aggregate.get('date_epoch', now.timestamp())
            record.pop('notification_from', '')
            if aggregate.get('ttl', -1) < 0:
                record['ttl'] = aggregate.get('ttl', -1)
            comment = {}
            comment['record_uid'] = aggregate['uid']
            comment['date'] = now.astimezone().isoformat()
            comment['auto'] = True
            if record_state == 'close':
                self.core.stats.inc('alert_closed', {'name': aggrule_name})
                if record.get('state') != 'close':
                    log.debug("OK received, closing alert")
                    aggregate_severity = aggregate.get('severity', 'unknown')
                    record_severity = record.get('severity', 'unknown')
                    comment['message'] = f"Auto closed: Severity {aggregate_severity} => {record_severity}"
                    comment['type'] = 'close'
                    record['state'] = 'close'
                    self.db.write('comment', comment)
                    record['comment_count'] = aggregate.get('comment_count', 0) + 1
                    return record
                else:
                    log.debug("OK received but the alert is already closed, discarding")
                    raise AbortAndUpdate(record)
            watched_fields = []
            for watched_field in watch:
                aggregate_field = dig(aggregate, *watched_field.split('.'))
                record_field = dig(record, *watched_field.split('.'))
                log.debug("Watched field %s: compare %s and %s", watched_field, record_field, aggregate_field)
                if record_field != aggregate_field:
                    watched_fields.append({'name': watched_field, 'old': aggregate_field, 'new': record_field})
            if watched_fields:
                log.debug("Alert %s Found updated fields from watchlist: %s", str(record['hash']), watched_fields)
                append_txt = []
                for watch_field in watched_fields:
                    append_txt.append("{watch_field['name']} ({watch_field['old']} => {watch_field['new']})")
                if record.get('state') == 'close':
                    comment['message'] = 'Auto re-opened from watchlist: {}'.format(', '.join(append_txt))
                    comment['type'] = 'open'
                    record['state'] = 'open'
                elif record.get('state') == 'ack':
                    comment['message'] = 'Auto re-escalated from watchlist: {}'.format(', '.join(append_txt))
                    comment['type'] = 'esc'
                    record['state'] = 'esc'
                else:
                    comment['message'] = 'New escalation from watchlist: {}'.format(', '.join(append_txt))
                    comment['type'] = 'comment'
                self.db.write('comment', comment)
                record['flapping_countdown'] = aggregate.get('flapping_countdown', flapping) - 1
            elif record.get('state') == 'close':
                comment['message'] = 'Auto re-opened'
                comment['type'] = 'open'
                record['state'] = 'open'
                self.db.write('comment', comment)
                record['flapping_countdown'] = aggregate.get('flapping_countdown', flapping) - 1
            elif (throttle < 0) or (now.timestamp() - aggregate.get('date_epoch', 0) < throttle):
                log.debug("Alert %s Time within throttle %s range, discarding", record['hash'], throttle)
                self.core.stats.inc('alert_throttled', {'name': aggrule_name})
                raise AbortAndUpdate(record)
            else:
                if record.get('state') == 'ack':
                    comment['type'] = 'esc'
                    record['state'] = 'esc'
                else:
                    comment['type'] = 'comment'
                comment['message'] = 'New escalation'
                self.db.write('comment', comment)
                record.pop('flapping_countdown', '')
            record['comment_count'] = aggregate.get('comment_count', 0) + 1
            if record.get('flapping_countdown', 0) < 0:
                log.debug("Alert %s is flapping, discarding", record['hash'])
                raise AbortAndUpdate(record)
        else:
            log.debug("Not found, creating a new aggregate")
            record['duplicates'] = 1
        record.pop('snoozed', '')
        record.pop('notifications', '')
        return record

    def reload_data(self):
        super().reload_data()
        aggregate_rules = []
        for aggrule in (self.data or []):
            aggregate_rules.append(AggregateruleObject(aggrule))
        self.aggregate_rules = aggregate_rules

