#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Contain the base classes for plugins'''

#from .basic.plugin import Plugin, Abort, AbortAndWrite, AbortAndUpdate

import json
from abc import ABC, abstractmethod
from typing import *
from uuid import UUID

import falcon
from falcon import Request, Response
from pydantic import BaseModel, ValidationError, Field, validator
from pydantic.generics import GenericModel

from snooze.database import Database, Pagination
from snooze.utils.functions import authorize
from snooze.utils.condition import guess_condition, Condition, AlwaysTrue, And
from snooze.utils.parser import parser
from snooze.utils.model import Partial, Record, ApiModel
from snooze.utils.typing import AuthPayload, SnoozeUser

class ProcessPlugin(ABC):
    '''Base class for process plugins'''
    @abstractmethod
    def process(self, record: Record) -> Record:
        '''Process a record, and return the processed record'''

class ActionPlugin(ABC):
    '''Base class for action plugins'''
    @abstractmethod
    def send(self, record: Record, content: str):
        '''Method called when the action is triggered'''

class SearchRequest(BaseModel):
    '''Represent the payload used for search requests.
    Supports 2 types of search format:
    * Query language: To interpret the query language in the search bar
    * Condition: Programmatic way to search
    If both fields are provided, they are combined with a `AND` (which happens
    when a user search is combined with a tab/environment filter for instance).
    If none of the field is provided, the search will be equivalent to `ALWAYS_TRUE`
    '''
    query_language: Optional[str] = Field(
        title='Query language',
        default=None,
        description='A query language representing a condition',
    )
    condition: Optional[Condition] = Field(
        default=None,
        description='A condition to search',
    )

    @validator('condition')
    def parse_json(cls, value):
        '''Parse a condition in json format'''
        if isinstance(value, str):
            value = json.loads(value)
        return value

    def to_condition(self) -> Condition:
        '''Return the condition that the search request represent'''
        if not self.query_language and not self.condition:
            return AlwaysTrue()
        conditions = []
        if self.query_language:
            conditions.append(parser(self.query_language))
        if self.condition:
            conditions.append(self.condition)
        if len(conditions) > 1:
            return And(conditions=conditions)
        else:
            return conditions[0]

def parse_uid(uid_str: str) -> UUID:
    '''Parse a UUID, and fails with 400 if invalid'''
    try:
        return UUID(uid_str)
    except (ValueError, TypeError) as err:
        raise falcon.HTTPBadRequest(description=f"UUID '{uid_str}' is invalid: {err}")

def parse_uids(uids_str: str) -> List[UUID]:
    '''Parse a list of UIDs, separated by '+', return 400 if invalid'''
    uids = []
    for uid_str in uids_str.split('+'):
        try:
            uid = UUID(uid_str)
        except (ValueError, TypeError) as err:
            raise falcon.HTTPBadRequest(description=f"Error parsing UIDs '{uids_str}': '{uid_str}' is invalid: {err}")
        uids.append(uid)
    return uids

Model = TypeVar('Model', bound=BaseModel)

def parse_model(data: dict, model: Type[Model]) -> Model:
    '''Parse a given BaseModel, and fail with 400 if invalid'''
    try:
        return model(**data)
    except ValidationError as err:
        raise falcon.HTTPBadRequest(description=f"Invalid {model}: {err}")

# TODO: Maybe can be useful?
class BasePluginRoute(ABC):
    def __init__(self,
        api: 'Api',
        model: Type[Type[ApiModel]],
        plugin: Optional[ApiPlugin] = None,
    ):
        self.endpoint = api.database[model]
        self.model = model
        self.plugin = plugin

class PluginMiddleware(Generic[Model], ABC):
    def __init__(self, database, plugin):
        self.database = database
        self.plugin = plugin

    @abstractmethod
    def on_create(self, document: Model):
        ...

    @abstractmethod
    def on_patch(self, search, partial: Partial[Model]):
        ...

    @abstractmethod
    def on_search(self, condition: Condition, pagination: Pagination):
        ...

class BaseApiRoute:
    '''Base route for all API managed objects'''
    endpoint: 'Endpoint'

    def __init__(self,
        api: 'Api',
        model: Type[ApiModel],
        plugin: Optional[ApiPlugin] = None,
        middlewares: List[PluginMiddleware] = None,
    ):
        self.endpoint = api.database[model]
        self.model = model
        self.plugin = plugin
        self.middlewares = middlewares

    @staticmethod
    def inject_auth_payload(auth: AuthPayload, media: dict):
        '''Use the authentication payload to add information to the input data'''
        media['snooze_user'] = SnoozeUser(**auth.dict())

    @authorize
    def on_get(self, req: Request, resp: Response):
        '''Return a list of Api objects (matching the search filter if provided)'''
        search_request = parse_model(req.params, SearchRequest)
        pagination = parse_model(req.params, Pagination)
        condition = search_request.to_condition()
        resp.media = self.endpoint.search(condition, pagination)
        resp.status = falcon.HTTP_OK

    @authorize
    def on_get_uid(self, _req: Request, resp: Response, uid_str: str):
        '''Fetch one document by UID'''
        uid = parse_uid(uid_str)
        document = self.endpoint.get_by_uid(UUID(uid))
        resp.media = document.dict()
        resp.status = falcon.HTTP_OK

    @authorize
    def on_post(self, req: Request, resp: Response):
        '''Insert a new object in the database. Will fail if the object exist
        or there is an object with the same primary fields'''
        self.inject_auth_payload(req.context.auth, req.media)
        document = parse_model(req.media, self.model)
        self.endpoint.create(document)
        resp.status = falcon.HTTP_CREATED
        resp.media = document.dict()

    @authorize
    def on_put_uid(self, req: Request, resp: Response, uid_str: str):
        '''Replace an object at a given UID'''
        uid = parse_uid(uid_str)
        self.inject_auth_payload(req.context.auth, req.media)
        document = parse_model(req.media, self.model)
        self.endpoint.replace(uid, document)
        resp.status = falcon.HTTP_OK
        resp.media = document.dict()

    @authorize
    def on_patch_uids(self, req: Request, resp: Response, uids_str: str):
        '''Update an object with a partial update (merge)'''
        uids = parse_uids(uids_str)
        self.inject_auth_payload(req.context.auth, req.media)
        patch = parse_model(req.media, Partial[self.model])
        resp.media = self.endpoint.patch_many(uids, patch)
        resp.status = falcon.HTTP_OK

    @authorize
    def on_delete_uids(self, _req: Request, resp: Response, uids_str: str):
        '''Correspond to delete operations'''
        uids = parse_uids(uids_str)
        resp.media = self.endpoint.delete_many(uids)
        resp.status = falcon.HTTP_OK

class ApiPlugin(ABC):
    '''Base class for a plugin managing an API object'''
    model: ClassVar[Optional[Type[ApiModel]]]
    route: ClassVar[Type[Any]]

    name: str
    database: Database

    def __init_subclass__(cls,
        model: Optional[Type[ApiModel]] = None,
        route: Type[Any] = BaseApiRoute,
    ):
        '''Classes inheriting this abstract class can pass arguments when inheriting,
        in order to set the class variables.
        Example:
        class CustomPlugin(ApiPlugin, model=CustomObject, route=CustomRoute):
            ...
        '''
        cls.model = model
        cls.route = route

    def __init__(self,
        core: 'Core',
    ):
        self.name = self.__class__.__name__.lower()
        self.database = core.database

    def load_routes(self, api: falcon.App):
        '''Add the routes related to the plugin to the falcon WSGI'''
        if self.model is not None:
            route_instance = self.route(api, self.model, plugin=self)
            api.add_route(f"/api/{self.name}", route_instance)
            api.add_route(f"/api/{self.name}" + '/{uid_str}', route_instance, suffix='uid')
            api.add_route(f"/api/{self.name}" + '/{uids_str}', route_instance, suffix='uids')

    def reload(self):
        '''Reload the plugin (data, setup, config, etc)'''
