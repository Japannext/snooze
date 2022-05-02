#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Module for referencing custom exceptions used by snooze'''

import traceback
import sys
from typing import Optional
from logging import getLogger

from falcon import Request

from snooze.utils.typing import Record

log = getLogger('snooze')

class NonResolvableHost(RuntimeError):
    '''Thrown when one member of the cluster address cnanot be resolved
    by DNS.'''
    def __init__(self, host: str):
        self.host = host
        super().__init__(f"DNS cannot resolve {host}")

class SelfNotInCluster(RuntimeError):
    '''Thrown when the running application addresses are not defined in the cluster
    configuration'''

class SelfTooMuchInCluster(RuntimeError):
    '''Thrown when the current node has too many entries of his addresses in the cluster
    configuration'''

class ProcessError(RuntimeError):
    '''Raised during the process of a record'''
    def __init__(self, err: Exception, record: Record):
        self.record = record
        record['error'] = str(err)
        record['traceback'] = traceback.format_exception(*sys.exc_info())
        super().__init__(str(err))

class ApiError(RuntimeError):
    '''Raised when the API failed to process a request'''
    def __init__(self, req: Request ,err: Exception, obj: dict, message: Optional[str] = None):
        obj['error'] = f"Error during validation: {err}"
        obj['traceback'] = traceback.format_exception(*sys.exc_info())
        self.path = req.path
        self.params = req.params
        self.rejected = obj
        super().__init__(f"API error: {message or str(err)}")

class ConditionValidationError(ApiError):
    '''Raised when the validation fails during an API call'''
    def __init__(self, req: Request, err: Exception, obj: dict):
        super().__init__(req, err, obj, f"Validation error: {err}")

class DatabaseError(RuntimeError):
    '''Wrapper for database errors (putting more info about each query)'''
    def __init__(self, operation: str, details: dict, err: Exception):
        self.operation = operation
        self.details = details
        super().__init__(self, f"Database error during {operation} ({details}): {err}")

class ImmutableFieldError(RuntimeError):
    '''Raised when an immutable field is being updated'''
    def __init__(self, collection: str, field: str):
        super().__init__(self, f"Field '{field}' in collection '{collection}' is immutable and cannot be updated")

