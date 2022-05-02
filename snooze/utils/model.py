#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''A series of classes representing the data model in the database'''

from datetime import datetime, timedelta
from enum import Enum
from typing import List, Any, Dict, Type, ClassVar, Optional, Generic, TypeVar
from uuid import UUID, uuid4
from logging import getLogger

import dateutil.parser
from pydantic import BaseModel, Extra, Field, root_validator, validator
from pydantic.generics import GenericModel

from snooze.utils.condition import Condition, AlwaysTrue
from snooze.utils.modification import Modification
from snooze.utils.typing import SnoozeUser

log = getLogger('snooze.utils.typing')

class RecordState(Enum):
    '''An enum representing the possible states of an alert'''
    ACK = 'ack'
    ESC = 'esc'
    CLOSED = 'closed'
    SHELVED = 'shelved'
    ERROR = 'error'
    EMPTY = ''

class MongodbMetadata(BaseModel):
    '''A base model for the MongodbMetadata(BaseMongodb) configuration class declaration'''
    collection: str = Field(
        description='The mongodb collection where the object will be stored',
    )
    primaries: set = Field(
        default_factory=set,
        description='A set of fields that should be considered as primary key (ensure uniqueness)'
    )
    immutables: set = Field(
        default_factory=set,
        description='A set of fields which should be considered immutables. If a user try to update on of'
        ' the given fields, an exception will be raised',
    )

class DatabaseEntry(BaseModel):
    '''Represent a document in the database'''
    uid: UUID = Field(default_factory=uuid4)
    date_epoch: float = Field(default_factory=lambda: datetime.now().timestamp())

    _mongodb: ClassVar[MongodbMetadata] = Field(
        title='MongodbMetadata(BaseMongodb) specific class configuration',
        description='Used by each new object to describe how it should be stored in mongodb',
    )

class UserEntry(DatabaseEntry):
    '''Represent a document created by a human or API user, as opposed to
    snooze internals and input plugins'''
    snooze_user: SnoozeUser = Field(
        title='Snooze user',
        description='User that last modified the object',
    )

class Record(DatabaseEntry, extra=Extra.allow):
    '''An alert entering the system'''
    _mongodb = MongodbMetadata(collection='record')

    timestamp: datetime = Field(default_factory=datetime.now)
    source: str = 'unknown'
    host: str = ''
    message: str = ''
    process: str = ''
    severity: str = 'unknown'
    environment: str = 'unknown'
    ttl: Optional[int] = None
    state: RecordState = RecordState.EMPTY


    @validator('timestamp')
    def dateutil_parser(cls, value):
        '''Parse the timestamp using dateutil.parser.'''
        if isinstance(value, str):
            try:
                return dateutil.parser.parse(value)
            except dateutil.parser.ParserError as err:
                log.warning('timestamp parsing error', exc_info=err)
                return datetime.now()
        return value

class Rule(UserEntry):
    _mongodb = MongodbMetadata(collection='rule')

    name: str = ''
    condition: Condition = Field(default_factory=AlwaysTrue, discriminator='type')
    modifications: List[Modification] = Field(default_factory=list)
    comment: str = ''


class AggregateRule(UserEntry):
    _mongodb = MongodbMetadata(collection='aggregaterule')

    name: str = ''
    fields : List[str] = Field(default_factory=list)
    watch : List[str] = Field(default_factory=list)
    throttle: timedelta = timedelta(minutes=15)
    comment: str = ''
    condition: Condition = Field(default_factory=AlwaysTrue, discriminator='type')


class SnoozeFilter(UserEntry):
    _mongodb = MongodbMetadata(collection='snooze')

    enabled: bool = True
    name: str = ''
    condition: Condition = Field(default_factory=AlwaysTrue, discriminator='type')
    #time_constraints: MultiConstraint = Field()
    comment: str = ''
    discard: bool = False
    # Computed
    sort: str


    @root_validator
    def compute_sort(cls, values: dict) -> Dict: # pylint: disable=no-self-argument,no-self-use
        '''Computing `sort`'''
        if 'sort' not in values:
            values['sort'] = values['time_constraint'].get_sort()

class NotificationFrequency(BaseModel):
    '''Parameters related to the frequency of sending notification in case of error.'''
    total: int = Field(
        default=1,
        description='Total number of retries to perform if the notification action fails',
    )
    delay: timedelta = Field(
        default=timedelta(seconds=0),
        description='Delay (in seconds) to wait before sending the notification. This will'
        ' deffer the notification in a dedicated queue. If the notification is acknowledged '
        'or closed before the delay is over, it will be removed from the deffered queue and'
        ' will not be sent.',
    )
    every: timedelta = Field(
        default=timedelta(seconds=0),
        description='Interval (in seconds) between retries if the notification fails',
    )

class Notification(UserEntry):
    _mongodb = MongodbMetadata(collection='notification')

    '''A rule that define what should happen when a record is not snoozed (and should be notified)'''
    enabled: bool = Field(
        default=True,
        description='Whether the notification is enabled or disabled',
    )
    name: str = Field(
        default='',
        description='Name of the notification',
    )
    condition: Condition = Field(
        discriminator='type',
        default_factory=AlwaysTrue,
        description='Condition for which the notification will be triggered.',
    )
    actions: List[str] = Field(
        default_factory=list,
        description='List of actions that should be triggered when the notification is triggered',
    )
    frequency: NotificationFrequency = Field(
        default_factory=NotificationFrequency,
    )
    comment: str = Field(
        default='',
        description='Comment associated with the notification',
    )

class Action(UserEntry):
    _mongodb = MongodbMetadata(collection='action', primaries={'name'})

    name: str = ''
    action: Any # TODO
    comment: str = ''
    # Computed
    pprint: str
    icon: str

class Widget(UserEntry):
    _mongodb = MongodbMetadata(collection='widget')

    enabled: bool = True
    name: str = ''
    widget: Any # TODO
    comment: str = ''
    # Computed
    pprint: str
    icon: str
    vue_component: str

class User(UserEntry):
    '''Represent a user data in the database'''
    _mongodb = MongodbMetadata(collection='user', primaries={'user', 'method'})

    name: str
    method: str
    roles: List[str] = Field(default_factory=list)
    groups: List[str] = Field(default_factory=list)
    last_login: datetime = Field(default_factory=datetime.now)

T = TypeVar('T', bound='BaseModel')

class Partial(Generic[T]):
    '''Partial[<Type>] returns a pydantic BaseModel identic to the given one,
    except all arguments are optional and defaults to None.
    This is intended to be used with partial updates.'''
    _types = {}

    def __class_getitem__(cls: Type[T], item: Type[Any]) -> Type[Any]:
        if isinstance(item, TypeVar):
            # Handle the case when Partial[T] is being used in type hints,
            # but T is a TypeVar, and the type hint is a generic. In this case,
            # the actual value doesn't matter.
            return item
        if item in cls._types:
            # If the value was already requested, return the same class. The main
            # reason for doing this is ensuring isinstance(obj, Partial[MyModel])
            # works properly, and all invocation of Partial[MyModel] return objects
            # that have the same class.
            new_model = cls._types[item]
        else:
            class new_model(item):
                '''Wrapper class to inherit the given class'''
            for _, field in new_model.__fields__.items():
                field.required = False
                if getattr(field, 'default_factory'):
                    field.default_factory = lambda: None
                else:
                    field.default = None
            cls._types[item] = new_model
        return new_model
