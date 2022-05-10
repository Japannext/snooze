#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#
'''
Objects and utils for representing conditions.
'''

from __future__ import annotations

import re
from abc import abstractmethod, ABC
from logging import getLogger
from typing import Any, Optional, List, ClassVar, Literal, Union, ForwardRef, Pattern

from pydantic import BaseModel, Field

from snooze.utils.functions import dig

log = getLogger('snooze.condition')

SCALARS = (str, int, float)

# Functions
def search(record: dict, field: str) -> Any:
    '''Searching a record while supporting . in field'''
    return dig(record, *field.split('.'))

def unsugar_regex(regex: str) -> str:
    '''Remove the leading and ending `/` if they are both present'''
    if len(regex) > 0 and regex[0] == '/' and regex[-1] == '/':
        regex = regex[1:-1]
    return regex

def lazy_search(value, string):
    '''Attempt to regex search a word in a string'''
    try:
        return re.search(str(value), str(string), flags=re.IGNORECASE)
    except TypeError:
        return False

# Exceptions
class OperationNotSupported(Exception):
    '''Exception raised when the condition requested doesn't exist'''
    def __init__(self, name):
        message = f"Condition '{name}' is not supported"
        super().__init__(message)

class ConditionInvalid(RuntimeError):
    '''Exception raised when there was an error when creating a condition,
    usually due to invalid inputs incompatible with the condition type'''
    def __init__(self, name, args, err):
        message = f"Error in condition '{name}' ({args}): {err}"
        super().__init__(message)

OperatorType = Literal[
    'AND', 'OR', 'NOT', # Logic
    '=', '!=', '>', '=>', '<', '<=', # Basic operators
    'ALWAYS_TRUE', 'EXISTS', 'SEARCH',
    'IN', 'CONTAINS', # Array based operators
]

# ForwardRefs
# https://pydantic-docs.helpmanual.io/usage/postponed_annotations/
And = ForwardRef('And')
Or = ForwardRef('Or')
Not = ForwardRef('Not')

# Classes
class ConditionBase(BaseModel, ABC):
    '''An abstract class for all conditions'''
    type: OperatorType = Field(..., include=True)

    @abstractmethod
    def match(self, record):
        '''Return true if the record match the condition'''

    #@abstractmethod
    def mongo_search(self) -> dict:
        '''Return the corresponding mongodb search'''
        raise NotImplementedError(f"Mongodb search is not implemented for {self.type}")

    def __and__(self, other: Condition) -> And:
        return And(conditions=[self, other])

    def __or__(self, other: Condition) -> Or:
        return Or(conditions=[self, other])

    def __invert__(self) -> Not:
        return Not(condition=self)

class BinaryOperator(ConditionBase, ABC):
    '''An abstract class to wrap binary operators'''
    display_name: ClassVar[Optional[str]] = None
    field: str
    value: Any

    def __str__(self):
        op_name = self.display_name or self.type.lower()
        return f"{self.field} {op_name} {repr(self.value)}"

class AlwaysTrue(ConditionBase):
    '''A condition that always return True for matching'''
    type: Literal['ALWAYS_TRUE'] = 'ALWAYS_TRUE'

    def match(self, record):
        return True
    def __str__(self):
        return '()'
    def mongo_search(self):
        return {}

# Logic
class Not(ConditionBase):
    '''Match the opposite of a given condition'''
    type: Literal['NOT'] = 'NOT'
    condition: Condition = Field(..., discriminator='type')

    def match(self, record):
        return not self.condition.match(record)

    def __str__(self):
        return f"!({self.condition})"

    def mongo_search(self):
        return {'$nor': self.condition.mongo_search()}

class And(ConditionBase):
    '''Match only if all conditions given in arguments match'''
    type: Literal['AND'] = 'AND'
    conditions: List[Condition]

    def match(self, record):
        return all(condition.match(record) for condition in self.conditions)
    def __str__(self):
        return '(' + ' & '.join(map(str, self.conditions)) + ')'
    def mongo_search(self):
        return {'$and': [condition.mongo_search() for condition in self.conditions]}

class Or(ConditionBase):
    '''Match only if one of the condition given in arguments match'''
    type: Literal['OR'] = 'OR'
    conditions: List[Condition]

    def match(self, record):
        return any(condition.match(record) for condition in self.conditions)
    def __str__(self):
        return '(' + ' | '.join(map(str, self.conditions)) + ')'
    def mongo_search(self):
        return {'$or': [condition.mongo_search() for condition in self.conditions]}

# Basic operations
class Equals(BinaryOperator):
    '''Match if the field of a record is exactly equal to a given value'''
    type: Literal['='] = '='

    def match(self, record):
        record_value = search(record, self.field)
        return record_value == self.value
    def mongo_search(self):
        return {self.field: self.value}

class NotEquals(BinaryOperator):
    '''Match if a field of a record is not equal to a given value'''
    type: Literal['!='] = '!='

    def match(self, record):
        record_value = search(record, self.field)
        return record_value != self.value
    def mongo_search(self):
        return {self.field: {'$ne', self.value}}

class GreaterThan(BinaryOperator):
    '''Match if the field of a record is strictly greater than a value.'''
    type: Literal['>'] = '>'

    def match(self, record):
        try:
            record_value = search(record, self.field)
            return record_value > self.value
        except TypeError as err: # Cannot be compared
            log.warning("%s > %s", repr(record_value), repr(self.value), exc_info=err)
            return False
    def mongo_search(self):
        return {self.field: {'$gt': self.value}}

class LowerThan(BinaryOperator):
    '''Match if the field of a record is strictly lower than a value.'''
    type: Literal['<'] = '<'

    def match(self, record):
        try:
            record_value = search(record, self.field)
            return record_value < self.value
        except TypeError as err: # Cannot be compared
            log.warning("%s < %s", repr(record_value), repr(self.value), exc_info=err)
            return False
    def mongo_search(self):
        return {self.field: {'$lt': self.value}}

class GreaterOrEquals(BinaryOperator):
    '''Match if the field of a record is greater than or equal to a value.'''
    type: Literal['>='] = '>='

    def match(self, record):
        try:
            record_value = search(record, self.field)
            return record_value >= self.value
        except TypeError as err: # Cannot be compared
            log.warning("%s >= %s", repr(record_value), repr(self.value), exc_info=err)
            return False
    def mongo_search(self):
        return {self.field: {'$gte': self.value}}

class LowerOrEquals(BinaryOperator):
    '''Match if the field of a record is lower than or equal a value.'''
    type: Literal['<='] = '<='

    def match(self, record):
        record_value = search(record, self.field)
        try:
            return record_value <= self.value
        except TypeError as err: # Cannot be compared
            log.warning("%s <= %s", repr(record_value), repr(self.value), exc_info=err)
            return False
    def mongo_search(self):
        return {self.field: {'$lte': self.value}}

# Complex operations
class Matches(ConditionBase):
    '''Match if the field of a record match a given regex. The regex can optionally
    start and end with `/`, to make it easier to spot in configuration. The regex method
    used is a search (`re.search`), so for strict matches, the `^` and `$` need to be used.
    '''
    type: Literal['MATCHES'] = 'MATCHES'
    field: str
    value: Pattern

    def __init__(self, **data):
        data['value'] = unsugar_regex(data['value'])
        ConditionBase.__init__(self, **data)

    def match(self, record):
        record_value = search(record, self.field)
        if record_value is None:
            return False
        try:
            return self.value.search(record_value) is not None
        except TypeError as err:
            log.warning("`record[%s] = %s` is not a string", self.field, repr(record_value), exc_info=err)
            return False
    def __str__(self):
        return f"{self.field} ~ /{self.value.pattern}/"

    def dict(self, **kwargs):
        '''Overriding for serializing the Pattern value'''
        data = ConditionBase.dict(self, **kwargs)
        if 'value' in data:
            data['value'] = data['value'].pattern
        return data
    def mongo_search(self):
        return {self.field: {'$regex': self.value.pattern, '$options': '-i'}}

class Exists(ConditionBase):
    '''Match if a given field exist and is not null in the record'''
    type: Literal['EXISTS'] = 'EXISTS'
    field: str

    def match(self, record):
        return search(record, self.field) is not None
    def __str__(self):
        return self.field + '?'
    def mongo_search(self):
        return {self.field: {'$exists': True}}

class Search(ConditionBase):
    '''Search a given string in the key/values of a record (stringify the record and
    search in the string)'''
    type: Literal['SEARCH'] = 'SEARCH'
    value: str

    def match(self, record):
        return self.value in str(record)
    def __str__(self):
        return f"(SEARCH {repr(self.value)})"

class ArrayContains(BinaryOperator):
    '''Match if it finds a given word/regex in a flatten list of object, or in a string'''
    type: Literal['CONTAINS'] = 'CONTAINS'

    def match(self, record):
        record_value = search(record, self.field)
        if isinstance(record_value, list):
            return self.value in record_value
        else:
            raise TypeError("") # TODO: Raise when record doesn't contain an array
            # we need to raise a warning

class InArray(ConditionBase):
    '''Match if a record field is in a given list of objects, or if
    the record field has any item matching a given condition.
    '''
    type: Literal['IN'] = 'IN'
    field: str
    value: List[Any]

    def match(self, record):
        record_value = search(record, self.field)
        return record_value in self.value

    """
    def __init__(self, args):
        super().__init__(args)
        self.field = args[2]
        self.value = args[1]
        if self.is_condition():
            self.mode = 'condition'
            self.condition = get_condition(self.value)
        else:
            self.mode = 'list'

    def is_condition(self):
        '''Detect if the provided argument is a condition or a scalar'''
        try:
            return self.value[0] in CONDITIONS
        except IndexError:
            return False

    def match(self, record):
        record_value = search(record, self.field)
        try:
            if self.mode == 'condition':
                return any(
                    self.condition.match(rec)
                    for rec in record_value
                )
            if self.mode == 'list':
                return any(
                    rec in flatten([self.value])
                    for rec in flatten([record_value])
                )
        except Exception as err:
            log.exception(err)
            return False
        # Unknown case
        log.warning("Unknown situation encountered for IN condition: condition=%s, record=%s",
            self._args, record)
        return False

    def __str__(self):
        if self.mode == 'condition':
            return f"({self.condition} in {self.field})"
        if self.mode == 'list':
            return f"({repr(self.value)} in {self.field})"
        return "???"
    """

Condition = Union[
    And,
    Or,
    Not,
    Exists,
    ArrayContains,
    InArray,
    Equals,
    NotEquals,
    Matches,
    GreaterOrEquals,
    LowerOrEquals,
    GreaterThan,
    LowerThan,
    Search,
    AlwaysTrue,
]

# ForwardRefs
# https://pydantic-docs.helpmanual.io/usage/postponed_annotations/
And.update_forward_refs()
Or.update_forward_refs()
Not.update_forward_refs()

class ConditionWrapper(BaseModel):
    '''A wrapper to be able to use pydantic defined constraints
    with the discriminator'''
    condition: Condition = Field(default_factory=AlwaysTrue, discriminator='type')

def guess_condition(data: dict) -> Condition:
    '''Return a condition given a dict representing the condition.
    Useful for testing.'''
    return ConditionWrapper(condition=data).condition
