'''Abstract classes for database'''

from abc import ABC, abstractmethod
from typing import Generic, TypeVar, List, Optional, Dict, Any, Type
from uuid import UUID

from pydantic import BaseModel

from snooze.utils.model import Partial
from snooze.utils.condition import Condition, AlwaysTrue
from snooze.utils.config import DatabaseConfig
from snooze.database.mongodb import MongodbDatabase

Model = TypeVar('Model', bound='BaseModel')

class Pagination(BaseModel):
    '''An object representing the pagination options of search'''
    order_by: str = '$natural'
    per_page: int = 10
    page_nb: int = 1
    asc: bool = True

    def to_skip(self):
        '''Return the number of entries to skip to access the desired content'''
        return self.per_page * (self.page_nb - 1)

class Database:
    '''Return the database of the correct type'''
    config: DatabaseConfig

    def __new__(cls, config: DatabaseConfig):
        if config.type == 'file':
            raise NotImplementedError()
        if config.type == 'mongo':
            return MongodbDatabase(config)
        raise NotImplementedError(f"Unknown database type '{config.type}'")

    @abstractmethod
    def __getitem__(self, model: Type[BaseModel]) -> Type['AbstractEndpoint']:
        ...

class DatabaseMiddleware(Generic[Model], ABC):
    '''A middleware to intercept requests to the database.
    Very useful to modify data before/after without having to
    deal with raw untyped data.'''
    def __init__(self, database, model):
        self.database = database
        self.model = model

    def __getattr__(self, method: str):
        if 

    @abstractmethod
    def on_create(self, document: Model):
        ...

    @abstractmethod
    def on_patch(self, search, partial: Partial[Model]):
        ...

    @abstractmethod
    def on_search(self, condition: Condition, pagination: Pagination):
        ...


class AbstractEndpoint(Generic[Model], ABC):
    '''A generic wrapper over a collection and a pydantic model'''

    @abstractmethod
    def search(self,
        condition: Condition = AlwaysTrue(),
        pagination: 'Pagination' = Pagination(),
    ) -> List[Model]:
        '''Search the collection with a condition'''

    @abstractmethod
    def count(self, condition: Condition = AlwaysTrue()) -> int:
        '''Return the total count of documents'''

    @abstractmethod
    def get_by_uid(self, uid: UUID) -> Model:
        '''Return a single object by uid'''

    @abstractmethod
    def get_by_index(self, index_dict: Dict[str, Any]) -> Model:
        '''Return an object by its indexed fields. Basically a search, but from a dict of indexes'''

    @abstractmethod
    def create(self, documents: List[Model]):
        '''Create a new document in the collection'''

    @abstractmethod
    def replace(self, uid: UUID, document: Model):
        '''Replace a document at a given UID with the given model'''

    @abstractmethod
    def patch_one(self, uid: UUID, document: Partial[Model]):
        '''Apply a partial update to a serie of document'''

    @abstractmethod
    def patch_many(self, uids: List[UUID], document: Partial[Model]):
        '''Apply a partial update to a serie of document'''

    @abstractmethod
    def delete_one(self, uid: UUID):
        '''Delete one document by uid'''

    @abstractmethod
    def delete_many(self, uids: List[UUID]):
        '''Delete documents by uid'''
