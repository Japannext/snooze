#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Typing utils for snooze'''

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
    port: str
    version: str
    healthy: bool

class HostPort(BaseModel):
    '''An object to represent a host-port pair'''
    host: str
    port: int = Field(5200)

class Pagination(TypedDict, total=False):
    '''A type hint for pagination options'''
    orderby: str
    nb_per_page: int
    page_nb: int
    asc: bool

class SettingUpdatePayload(BaseModel):
    filename: str
    conf: dict
    reload: bool

class Widget(BaseModel):
    '''A widget representation in the config'''
    widget_name: Optional[str] = None
    icon: str
    vue_component: str
    form: dict
