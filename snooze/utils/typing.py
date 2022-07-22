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

Condition = NewType('Condition', list)

Pagination = NewType('Pagination', dict)

DuplicatePolicy = Literal['insert', 'reject', 'replace', 'update']

# We're listing only the ones we might raise.
# It will not affect performance to list many, since falcon keep a dictionary of exception => handler.
HTTPUserErrors = (
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

class Widget(BaseModel):
    '''A widget representation in the config'''
    widget_name: Optional[str] = None
    icon: str
    vue_component: str
    form: dict
