#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Typing utils for snooze'''

from datetime import datetime
from typing import NewType, List, Literal, Optional, TypedDict, Union, Generic, TypeVar

import falcon
import jinja2
from pydantic import BaseModel, Field, ValidationError
from pydantic.fields import ModelField
from pydantic.generics import GenericModel

RecordUid = NewType('RecordUid', str)
Record = NewType('Record', dict)
Rule = NewType('Rule', dict)
AggregateRule = NewType('AggregateRule', dict)
SnoozeFilter = NewType('SnoozeFilter', dict)

Config = NewType('Config', dict)
Condition = NewType('Condition', list)
ConditionOrUid = Optional[Union[str, list]]


DuplicatePolicy = Literal['insert', 'reject', 'replace', 'update']

# We're listing only the ones we might raise.
# It will not affect performance to list many, since falcon keep a dictionary of exception => handler.
USER_ERRORS = (
    # HTTP 400
    falcon.HTTPBadRequest, falcon.HTTPInvalidHeader, falcon.HTTPInvalidParam, falcon.HTTPMissingParam,
    # HTTP 401
    falcon.HTTPUnauthorized,
    # HTTP 403
    falcon.HTTPForbidden,
    # HTTP 404
    falcon.HTTPNotFound, falcon.HTTPRouteNotFound,
)

class AuthorizationPolicy(BaseModel):
    '''A list of authorized policy for read and write'''
    read: List[str] = Field(default_factory=list)
    write: List[str] = Field(default_factory=list)

class RouteArgs(BaseModel):
    '''Description of the arguments passed to a route'''
    class_name: Optional[str] = Field('Route', alias='class')
    desc: Optional[str] = None
    primary: List[str] = Field(default_factory=list)
    duplicate_policy: DuplicatePolicy = 'update'
    authorization_policy: Optional[AuthorizationPolicy]
    check_permissions: bool = False
    check_constant: List[str] = Field(default_factory=list)
    inject_payload: bool = False
    prefix: str = '/api'

class PeerStatus(BaseModel):
    '''A dataclass containing the status of one peer'''
    host: str
    port: int
    version: str
    healthy: bool

class HostPort(BaseModel):
    '''An object to represent a host-port pair'''
    host: str = Field(
        required=True,
        description='The host address to reach (IP or resolvable hostname)',
    )
    port: int = Field(
        default=5200,
        description='The port where the host is expected to listen to'
    )

class Pagination(TypedDict, total=False):
    '''A type hint for pagination options'''
    orderby: str
    nb_per_page: int
    page_nb: int
    asc: bool

class Widget(BaseModel):
    '''A widget representation in the config'''
    widget_name: Optional[str] = None
    icon: str
    vue_component: str
    form: dict

class Query(BaseModel):
    ql: str
    field: str

class SnoozeUser(BaseModel):
    '''Represent the minimum information to uniquely identify a user in the system'''
    name: str = Field(
        title='User name',
        description='Name of the user in the system. Depends on the method.',
    )
    method: str = Field(
        title='Authentication method',
        description='Authentication method used by the user',
    )

class AuthPayload(BaseModel):
    '''An object representing the authentication payload that will
    be contained in the JWT token'''
    name: str
    method: str
    roles: List[str] = Field(default_factory=list)
    permissions: List[str] = Field(default_factory=list)
    groups: List[str] = Field(default_factory=list)

class ProfileGeneral(BaseModel):
    '''Represent a user profile in the database'''
    name: str
    method: str
    display_name: Optional[str] = None
    email: Optional[str] = None

class ProfilePreferences(BaseModel):
    '''Represent the preference page of a user'''
    name: str
    method: str
    default_page: str = '/record'
    theme: str = 'default'

T = TypeVar('T')

class JinjaTemplate(Generic[T]):
    '''A Jinja template type for pydantic (support for validation/deserialization)'''
    raw: T
    template: Optional[jinja2.Template]

    def __init__(self, raw: T):
        self.raw = raw
        if isinstance(raw, str) and '{' in raw:
            self.template = jinja2.Template(raw)
        else:
            self.template = None

    @classmethod
    def __get_validators__(cls):
        yield cls.validate

    @classmethod
    def __modify_schema__(cls, field_schema):
        pass

    @classmethod
    def validate(cls, value, field: ModelField):
        '''Validate the jinja template'''
        if isinstance(value, str):
            value = cls(value)
        if not field.sub_fields:
            return value
        inner_type = field.sub_fields[0]
        _, error = inner_type.validate(value.raw, {}, loc='raw')
        if error:
            raise ValidationError([error], cls)
        return value

    def resolve(self, env: dict) -> T:
        '''Render the template if it's a jinja template'''
        if self.template:
            return self.template.render(env)
        else:
            return self.raw

    def __repr__(self):
        return f"JinjaTemplate({repr(self.raw)})"
