#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''A module for managing time constraint objects, mainly used by the
snooze and notification core plugins'''

from __future__ import annotations

from abc import ABC, abstractmethod
from datetime import datetime, timedelta, time
from enum import Enum
from logging import getLogger
from typing import List, Optional, Tuple, Literal, Dict, Union, ForwardRef

import dateutil.parser
from pydantic import BaseModel, Field, validator

from snooze.utils.typing import Record

DatetimeRange = Tuple[datetime, datetime]

log = getLogger('snooze.time_constraints')

def get_record_date(record: Record) -> datetime:
    '''Extract the date of the record and return a `datetime` object'''
    if record.get('timestamp'):
        record_date = dateutil.parser.parse(record['timestamp']).astimezone()
    elif record.get('date_epoch'):
        record_date = datetime.fromtimestamp(record['date_epoch']).astimezone()
    else:
        record_date = datetime.now().astimezone()
    return record_date

OperatorType = Literal['datetime', 'time', 'weekdays', 'AND', 'OR', 'NOT']

# ForwardRefs
# https://pydantic-docs.helpmanual.io/usage/postponed_annotations/
And = ForwardRef('And')
Or = ForwardRef('Or')
Not = ForwardRef('Not')

class TemporalConstraintBase(BaseModel, ABC):
    '''A base class for time constraints'''
    type: OperatorType = Field(...)

    class Config:
        allow_population_by_field_name = True

    @abstractmethod
    def match(self, record_date: datetime) -> bool:
        '''Return if the temporal constraint match the provided record date'''

    def __and__(self, other: TemporalConstraint) -> And:
        return And(constraints=[self, other])

    def __or__(self, other: TemporalConstraint) -> Or:
        return Or(constraints=[self, other])

    def __invert__(self) -> Not:
        return Not(constraint=self)

    def __str__(self):
        return "TemporalConstraint"

# Logic
class Not(TemporalConstraintBase):
    '''Match the opposite of a given condition'''
    type: Literal['NOT'] = 'NOT'
    constraint: TemporalConstraint = Field(..., discriminator='type')

    def match(self, record_date: datetime):
        return not self.constraint.match(record_date)

    def __str__(self):
        return '!' + str(self.constraint)

    def mongo_search(self):
        return {'$nor': self.constraint.mongo_search()}

class And(TemporalConstraintBase):
    '''Match only if all constraints given in arguments match'''
    type: Literal['AND'] = 'AND'
    constraints: List[TemporalConstraint]

    def match(self, record_date: datetime):
        return all(constraint.match(record_date) for constraint in self.constraints)
    def __str__(self):
        return '(' + ' & '.join(map(str, self.constraints)) + ')'
    def mongo_search(self):
        return {'$and': [constraint.mongo_search() for constraint in self.constraints]}

class Or(TemporalConstraintBase):
    '''Match only if one of the constraint given in arguments match'''
    type: Literal['OR'] = 'OR'
    constraints: List[TemporalConstraint]

    def match(self, record_date: datetime):
        return any(constraint.match(record_date) for constraint in self.constraints)
    def __str__(self):
        return '(' + ' | '.join(map(str, self.constraints)) + ')'
    def mongo_search(self):
        return {'$or': [constraint.mongo_search() for constraint in self.constraints]}

class DatetimeConstraint(TemporalConstraintBase):
    '''A time constraint using fixed dates.
    Features:
        * Before a fixed date
        * After a fixed date
        * Between two fixed dates
    '''
    type: Literal['datetime'] = 'datetime'
    date_from: Optional[datetime] = Field(
        alias='from',
        default=None,
    )
    date_until: Optional[datetime] = Field(
        alias='until',
        default=None,
    )

    @validator('date_from', 'date_until', check_fields=False)
    def dateutil_parser(cls, value):
        '''Parse strings with dateutil.parser if a string is provided'''
        if isinstance(value, str):
            value = dateutil.parser.parse(value).astimezone()
        return value

    def match(self, record_date: datetime) -> bool:
        '''Perform a fixed date matching'''
        date_from = self.date_from
        date_until = self.date_until
        if date_from and date_until:
            return date_from <= record_date <= date_until
        elif (not date_from) and date_until:
            return record_date <= date_until
        elif date_from and (not date_until):
            return date_from <= record_date
        else:
            return False

    def __str__(self):
        return f"({self.date_from} -> {self.date_until})"

class WeekdayEnum(Enum):
    '''An enum for numeric weekdays'''
    SUNDAY    = 0
    MONDAY    = 1
    TUESDAY   = 2
    WEDNESDAY = 3
    THURSDAY  = 4
    FRIDAY    = 5
    SATURDAY  = 6

class WeekdaysConstraint(TemporalConstraintBase):
    '''A constraint on the days of the week
    Features:
        * Match certain days of the week
    '''
    type: Literal['weekdays'] = 'weekdays'
    week: Dict[WeekdayEnum, bool] = Field(default_factory=dict)

    def match(self, record_date: datetime) -> bool:
        weekday = WeekdayEnum(int(record_date.strftime('%w')))
        return self.week.get(weekday, False)
    def __str__(self):
        weekdays = [weekday.name for weekday, enabled in self.week.items() if enabled]
        return f"({' '.join(weekdays)})"

class TimeConstraint(TemporalConstraintBase):
    '''A time constraint that has a daily period.
    Features:
        * Match before/after/between fixed hours
        * Support hours over midnight (`from` lower than `until`)
    '''
    type: Literal['time'] = 'time'
    time1: Optional[time] = Field(
        alias='from',
        default=None,
    )
    time2: Optional[time] = Field(
        alias='until',
        default=None,
    )

    @validator('time1', 'time2', check_fields=False)
    def dateutil_parser(cls, value):
        '''Parse string times using dateutil.parser. Return the time with the timezone'''
        if isinstance(value, str):
            value = dateutil.parser.parse(value).astimezone().timetz()
        return value

    def get_intervals(self, record_date: datetime) -> List[DatetimeRange]:
        '''Return the an array of datetime intervals depending on the `from`
        and `until` time, and the date of the record. The intervals will all be
        ordered.'''
        day = timedelta(days=1)
        date1 = datetime.combine(record_date, self.time1)
        date2 = datetime.combine(record_date, self.time2)
        if date2 < date1:
            return [(date1 - day, date2), (date1, date2 + day)]
        return [(date1, date2)]

    def match(self, record_date: datetime) -> bool:
        '''Match a daily periodic time constraint'''
        rd = record_date.astimezone()
        if self.time1 and self.time2:
            intervals = self.get_intervals(rd.date())
            return any(date1 <= rd <= date2 for date1, date2 in intervals)
        elif self.time1 and not self.time2:
            date1 = datetime.combine(rd.date(), self.time1)
            return date1 <= rd
        elif self.time2 and not self.time1:
            date2 = datetime.combine(rd.date(), self.time2)
            return rd <= date2
        else:
            return True

    def __str__(self):
        return f"({self.time1} -> {self.time2})"

TemporalConstraint = Union[
    # Logic
    And, Or, Not,
    # Constraints
    DatetimeConstraint,
    WeekdaysConstraint,
    TimeConstraint,
]

# ForwardRefs
# https://pydantic-docs.helpmanual.io/usage/postponed_annotations/
And.update_forward_refs()
Or.update_forward_refs()
Not.update_forward_refs()

class ConstraintWrapper(BaseModel):
    '''A wrapper to be able to use pydantic defined constraints
    with the discriminator'''
    constraint: TemporalConstraint = Field(..., discriminator='type')

def guess_constraint(data: dict) -> TemporalConstraint:
    '''Return a constraint given a dict representing the temporal constraint.
    Useful for testing.'''
    return ConstraintWrapper(constraint=data).constraint
