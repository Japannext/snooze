#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Typing utils for snooze'''

from typing import NewType, List, Literal, Optional, List

from pydantic import BaseModel, Field

RecordUid = NewType('RecordUid', str)
Record = NewType('Record', dict)
Rule = NewType('Rule', dict)
AggregateRule = NewType('AggregateRule', dict)
SnoozeFilter = NewType('SnoozeFilter', dict)

Config = NewType('Config', dict)
Condition = NewType('Condition', list)

DuplicatePolicy = Literal['insert', 'reject', 'replace', 'update']

class AuthorizationPolicy(BaseModel):
    '''A list of authorized policy for read and write'''
    read: List[str]
    write: List[str]

class RouteArgs(BaseModel):
    '''Description of the arguments passed to a route'''
    class_name: Optional[str] = Field(alias='class')
    primary: Optional[str] = None
    duplicate_policy: DuplicatePolicy = 'update'
    authorization_policy: Optional[AuthorizationPolicy]
    check_permissions: bool = False
    check_constant: Optional[str] = None
    inject_payload: bool = False
    prefix: str = '/api'
