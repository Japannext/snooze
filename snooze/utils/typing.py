#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Typing utils for snooze'''

from datetime import datetime
from typing import NewType, List, Literal, Optional, TypedDict, Union

from pydantic import BaseModel, Field

RecordUid = NewType('RecordUid', str)
Record = NewType('Record', dict)
Rule = NewType('Rule', dict)
AggregateRule = NewType('AggregateRule', dict)
SnoozeFilter = NewType('SnoozeFilter', dict)

Config = NewType('Config', dict)
Condition = NewType('Condition', list)
ConditionOrUid = Optional[Union[str, list]]

DuplicatePolicy = Literal['insert', 'reject', 'replace', 'update']

class AuthorizationPolicy(BaseModel):
    '''A list of authorized policy for read and write'''
    read: Optional[List[str]]
    write: Optional[List[str]]

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

class AuthPayload(BaseModel):
    '''An object representing the authentication payload that will
    be contained in the JWT token'''
    username: str
    method: str
    roles: List[str] = Field(default_factory=list)
    permissions: List[str] = Field(default_factory=list)
    groups: List[str] = Field(default_factory=list)

class Profile(BaseModel):
    '''Represent a user and its preferences in the database'''
    username: str
    method: str
    display_name: Optional[str] = None
    email: Optional[str] = None
