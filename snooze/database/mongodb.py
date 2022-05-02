'''Mongodb wrapper for Pydantic models'''

from uuid import UUID
from typing import TypeVar, Type, List, Set, Generic, Generator, Literal, Dict, Any
from logging import getLogger
from contextlib import contextmanager
from functools import wraps
from datetime import timedelta

from pymongo import MongoClient
from pymongo.database import Database
from pymongo.collection import Collection
from pymongo.operations import IndexModel
from pymongo import ASCENDING, HASHED
from pydantic import BaseModel, Field, ValidationError
from bson.codec_options import CodecOptions, TypeCodec, TypeRegistry

from snooze.utils.condition import Condition, AlwaysTrue
from snooze.utils.model import DatabaseEntry, UserEntry, MongodbMetadata, Partial
from snooze.utils.exceptions import DatabaseError, ImmutableFieldError

log = getLogger('snooze.database.mongodb')

Model = TypeVar('Model', bound=DatabaseEntry)

class Pagination(BaseModel):
    '''An object representing the pagination options of search'''
    order_by: str = '$natural'
    per_page: int = 10
    page_nb: int = 1
    asc: bool = True

    def to_skip(self):
        '''Return the number of entries to skip to access the desired content'''
        return self.per_page * (self.page_nb - 1)

class SearchResult(Generic[Model], BaseModel):
    '''Represent a database search result'''
    total_count: int = Field(
        description='Total count of objects in the database, ignoring the pagination options',
    )
    data: List[Model] = Field(
        default_factory=list,
        description='Result of a database search operation'
    )

class BulkResult(Generic[Model], BaseModel):
    '''Represent a bulk result during an operation'''
    operation: Literal['create', 'update', 'delete']
    success: List[Model]
    failure: List[Model]

class Timedelta(TypeCodec):
    '''A codec to represent datetime.timedelta in BSON.
    Define how its conversion to python is done.'''
    python_type = timedelta
    bson_type = float

    def transform_python(self, value: float) -> timedelta:
        return timedelta(seconds=value)

    def transform_bson(self, value: timedelta) -> float:
        return value.total_seconds()

TYPE_REGISTRY = TypeRegistry([Timedelta()])
CODEC_OPTIONS = CodecOptions(type_registry=TYPE_REGISTRY, tz_aware=True, uuid_representation=4)

def wrap_exception(function):
    '''Wrap an Endpoint method exception so we get more information about the
    query that made it fail'''
    def wrapper(endpoint, *args, **kwargs):
        try:
            return function(endpoint, *args, **kwargs)
        except ImmutableFieldError as err:
            raise err
        except Exception as err:
            details = {'collection': endpoint.collection, 'args': args, 'kwargs': kwargs}
            raise DatabaseError(function.__name__, details, err) from err
    return wrapper

@contextmanager
def transaction(client: MongoClient, **kwargs):
    '''Wrapper context manager for MongoDB transactions'''
    with client.start_session() as session:
        with session.start_transaction(**kwargs):
            yield session

def get_collection(database: Database, collection_name: str):
    '''Ensure the existence, and return the collection of a database
    with all codec options.'''
    return database.get_collection(collection_name, codec_options=CODEC_OPTIONS)

class Endpoint(Generic[Model]):
    '''A generic wrapper over a collection and a pydantic model'''

    collection: Collection
    config: MongodbMetadata
    immutables: Set[str]

    def __init__(self, database: Database, model: Type[Model]):
        self.model = model
        self.config = model._mongodb
        self.client = database.client
        # Collection
        self.collection = get_collection(database, self.config.collection)
        # Indexes
        indexes = self._prepare_indexes(self.config.primaries)
        self._ensure_indexes(indexes)
        # Immutable fields
        self.immutables = self.config.immutables

    def _prepare_indexes(self, primaries: Set[str]) -> List[IndexModel]:
        '''Generate the objects representing the indexes to ensure'''
        indexes = []
        uid_index = IndexModel([('uid', HASHED)], name='uid', unique=True)
        indexes.append(uid_index)
        if self.config.primaries:
            index_keys = [(primary, ASCENDING) for primary in primaries]
            primary_index = IndexModel(index_keys, name='primary', unique=True)
            indexes.append(primary_index)
        return indexes

    def _ensure_indexes(self, indexes: List[IndexModel]):
        '''Ensure a given list of indexes exists'''
        # Using a transaction to avoid bootstrap conflicts at startup
        with transaction(self.client) as session:
            current_indexes = self.collection.index_information()
            for index in indexes:
                name = index.document['name']
                # Create index if not present
                if name not in current_indexes:
                    log.info("database[%s] creating index '%s'", self.config.collection, name)
                    self.collection.create_indexes(indexes=[index], session=session)
                    continue
                current_index = current_indexes[name]
                # Recreate index if keys don't match
                if set(current_index['key']) != set(index.document['key']):
                    log.warning("database[%s] index '%s' key change: %s => %s",
                        self.config.collection, name, current_index['key'], index.document['key'])
                    log.warning("database[%s] deleting index '%s'", self.config.collection, name)
                    self.collection.drop_index(name, session=session)
                    log.warning("database[%s] creating index '%s'", self.config.collection, name)
                    self.collection.create_index(session=session)
                    continue

    @wrap_exception
    def search(self,
        condition: Condition = AlwaysTrue(),
        pagination: Pagination = Pagination(),
    ) -> List[Model]:
        '''Search the collection with a condition'''
        mongo_search = condition.mongo_search()
        cursor = self.collection \
            .find(mongo_search) \
            .skip(pagination.to_skip()) \
            .limit(pagination.per_page) \
            .sort(pagination.order_by, pagination.asc)
        return list(self._fit_model(cursor))

    @wrap_exception
    def count(self, condition: Condition = AlwaysTrue()) -> int:
        '''Return the total count of documents'''
        mongo_search = condition.mongo_search()
        return self.collection.count_documents(mongo_search)

    @wrap_exception
    def get_by_uid(self, uid: UUID) -> Model:
        '''Return a single object by uid'''
        result = self.collection.find_one({'uid': uid})
        log.debug("get_by_uid(%s): %s", uid, result)
        return self.model(**result) if result else None

    @wrap_exception
    def get_by_index(self, index_dict: Dict[str, Any]) -> Model:
        '''Return an object by its indexed fields. Basically a search, but from a dict of indexes'''
        self._warn_foreign_indexes(index_dict.keys())
        result = self.collection.find_one(index_dict)
        return self.model(**result)

    def _warn_foreign_indexes(self, indexes: List[str]):
        '''Issue warning logs when the requested indexes are not among the defined primaries.
        The function will still work, but without the expected performance of an index.'''
        foreign_indexes = [index for index in indexes if index not in self.config.primaries]
        for index in foreign_indexes:
            log.warning("In database[%s].get_by_index(), %s not in primary indexes (%s)",
                self.config.collection, index, self.config.primaries)

    @wrap_exception
    def create(self, documents: List[Model]):
        '''Create a new document in the collection'''
        documents = [doc.dict() for doc in documents]
        log.debug('create(%s)', documents)
        result = self.collection.insert_many(documents)
        return result

    @wrap_exception
    def replace(self, uid: UUID, document: Model):
        '''Replace a document at a given UID with the given model'''
        document.uid = uid
        result = self.collection.find_one_and_replace({'uid': uid}, document.dict())
        return result

    @wrap_exception
    def patch_one(self, uid: UUID, document: Partial[Model]):
        '''Apply a partial update to a serie of document'''
        for immutable in self.immutables:
            if document[immutable] is not None:
                raise ImmutableFieldError(self.collection.name, immutable)
        result = self.collection.update_one({'uid': uid}, {'$set': document.dict(exclude_none=True)})
        return result

    @wrap_exception
    def patch_many(self, uids: List[UUID], document: Partial[Model]):
        '''Apply a partial update to a serie of document'''
        for immutable in self.immutables:
            if document[immutable] is not None:
                raise ImmutableFieldError(self.collection.name, immutable)
        result = self.collection.update_many({'uid': {'$in': uids}}, {'$set': document.dict(exclude_none=True)})
        return result

    @wrap_exception
    def delete_one(self, uid: UUID):
        '''Delete one document by uid'''
        result = self.collection.delete_one({'uid': uid})
        return result

    @wrap_exception
    def delete_many(self, uids: List[UUID]):
        '''Delete documents by uid'''
        result = self.collection.delete_many({'uid': {'$in': uids}})
        return result

    def _fit_model(self, cursor) -> Generator[Model, None, None]:
        '''A generator that take a cursor, and return a '''
        for document in cursor:
            try:
                yield self.model(**document)
            except ValidationError as err:
                uid = document.get('uid', 'unknown')
                log.warning("Invalid object in database[%s] at uid=%s",
                    self.config.collection, uid, exc_info=err)
                continue

    def _check_immutables(self, new_docs: List[Model]):
        '''Check a list of documents try to update an immutable field'''
        uids = [str(doc.uid) for doc in new_docs]
        cursor = self.collection.find({'uid': {'$in': uids}}, projection=list(self.immutables))
        old_docs = self._fit_model(cursor)
        docdict = {doc.uid: doc for doc in new_docs}
        for old_doc in old_docs:
            new_doc = docdict[old_doc.uid]
            for immutable in self.immutables:
                if old_doc[immutable] != new_doc[immutable]:
                    raise ImmutableFieldError(self.collection.name, immutable)
