#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Module for managing modification objects.
Modifications are used by the rule core plugin to modify records
automatically based on a rule.
'''

from abc import abstractmethod
from logging import getLogger
from typing import Literal, Union, List, Any, Pattern

import jinja2
from pydantic import BaseModel, Field

from snooze.utils.typing import Record, JinjaTemplate

log = getLogger('snooze.utils.modification')

class OperationNotSupported(Exception):
    '''Exception raised when the modification type doesn't exist'''
    def __init__(self, name):
        message = f"Modification type '{name}' is not supported"
        super().__init__(message)

class ModificationInvalid(RuntimeError):
    '''Exception raise when there was an error when creating a modification'''
    def __init__(self, name, args, err):
        message = f"Error in modification '{name}' ({args}): {err}"
        super().__init__(message)

OperationType = Literal[
    # Field operations
    'SET', 'DELETE',
    # Array operations
    'ARRAY_APPEND',
    'ARRAY_DELETE',
    # Regex operations
    'REGEX_PARSE', 'REGEX_SUB',
    # Others
    'KV_SET',
]

class ModificationBase(BaseModel):
    '''A class to represent a modification'''
    type: OperationType = Field(...)
    core: 'Optional[Core]' = None

    @abstractmethod
    def modify(self, record: Record) -> bool:
        '''Called when the modification should be applied to a record'''

    def resolve(self, record: Record, fields: List[str]):
        '''Resolve jinja templates for all fields using the record as a dict'''
        return [
            jinja2.Template(self[field]).render(record)
            for field in fields
        ]

    def __str__(self):
        return ''

class SetOperation(ModificationBase):
    '''Set a key to a given value'''
    type: Literal['SET'] = 'SET'
    field: JinjaTemplate[str]
    value: JinjaTemplate[Any]

    def modify(self, record: Record) -> bool:
        field = self.field.resolve(record)
        value = self.value.resolve(record)
        try:
            return_code = bool(value and record.get(field) != value)
            record[field] = value
            return return_code
        except Exception:
            return False

    def __str__(self):
        return f"record[{self.field}] = {repr(self.value)}"

class DeleteOperation(ModificationBase):
    '''Delete a given key'''
    type: Literal['SET'] = 'SET'
    field: JinjaTemplate[str]

    def modify(self, record: Record) -> bool:
        field = self.field.resolve(record)
        try:
            del record[field]
            return True
        except KeyError:
            return False

    def __str__(self):
        return f"del record[{self.field}]"

class ArrayAppendOperation(ModificationBase):
    '''Append an element to a key, if this key is an array/list'''
    type: Literal['ARRAY_APPEND'] = 'ARRAY_APPEND'
    field: JinjaTemplate[str]
    value: JinjaTemplate[Any]

    def modify(self, record: Record) -> bool:
        field = self.field.resolve(record)
        value = self.value.resolve(record)
        array = record.get(field)
        if isinstance(array, list):
            array += value
            return True
        else:
            return False

    def __str__(self):
        return f"record[{self.field}].append({repr(self.value)})"

class ArrayDeleteOperation(ModificationBase):
    '''Delete an element from an array/list, by value'''
    type: Literal['ARRAY_APPEND'] = 'ARRAY_APPEND'
    field: JinjaTemplate[str]
    value: JinjaTemplate[Any]

    def modify(self, record: Record) -> bool:
        field = self.field.resolve(record)
        value = self.value.resolve(record)
        try:
            record[field].remove(value)
            return True
        except (ValueError, KeyError):
            return False

    def __str__(self):
        return f"record[{self.field}].remove({repr(self.value)})"

class RegexParse(ModificationBase):
    '''Given a key and a regex with named capture groups, parse the
    key's value, and merge the captured elements with the record'''
    type: Literal['REGEX_PARSE'] = 'REGEX_PARSE'
    field: JinjaTemplate[str]
    regex: Pattern

    def modify(self, record: Record) -> bool:
        try:
            field = self.field.resolve(record)
            results = self.regex.search(record[field])
            if results:
                for name, value in results.groupdict({}).items():
                    record[name] = value
                return True
            return False
        except KeyError:
            return False

class RegexSub(ModificationBase):
    '''Apply a regex search and replace expression to a key's value'''
    type: Literal['REGEX_SUB'] = 'REGEX_SUB'
    field: JinjaTemplate[str]
    regex: Pattern
    sub: JinjaTemplate[str]
    out_field: str

    def modify(self, record: Record) -> bool:
        field = self.field.resolve(record)
        sub = self.sub.resolve(record)
        try:
            new_value = self.regex.sub(sub, record[field])
            record[self.out_field] = new_value
            return True
        except (KeyError, TypeError):
            return False

class KvSet(ModificationBase):
    '''Match the key's value with the corresponding value from the kv core plugin'''
    dictionary: JinjaTemplate[str]
    field: JinjaTemplate[str]
    out_field: JinjaTemplate[str]

    def modify(self, record: Record) -> bool:
        try:
            field = self.field.resolve(record)
            dictionary = self.dictionary.resolve(record)
            out_field = self.out_field.resolve(record)
            out_value = self.core.get_core_plugin('kv').get(dictionary, record.get(field))
            log.debug("Found key-value: %s[%s] = %s", dictionary, record.get(field), out_value)
            record[out_field] = out_value
            return True
        except (KeyError, IndexError):
            return False

Modification = Union[
    SetOperation,
    DeleteOperation,
    ArrayAppendOperation,
    ArrayDeleteOperation,
    RegexParse,
    RegexSub,
    KvSet,
]

def guess_modification(data: dict) -> Modification:
    ...
