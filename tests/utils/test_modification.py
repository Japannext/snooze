#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

from logging import getLogger

import pytest
from pydantic import ValidationError

from snooze.utils.modification import *

log = getLogger('tests')

def test_modification_set():
    record = {'a': 1, 'b': 2}
    modification = SetOperation(field='c', value=3)
    return_code = modification.modify(record)
    assert record == {'a': 1, 'b': 2, 'c': 3}
    assert return_code

def test_modification_delete():
    record = {'a': 1, 'b': 2}
    modification = DeleteOperation(field='b')
    return_code = modification.modify(record)
    assert record == {'a': 1}
    assert return_code

def test_modification_array_append():
    record = {'a': 1, 'b': ['1', '2', '3']}
    modification = ArrayAppendOperation(field='b', value='4')
    return_code = modification.modify(record)
    assert record == {'a': 1, 'b': ['1', '2', '3', '4']}
    assert return_code

def test_modification_array_delete():
    record = {'a': 1, 'b': ['1', '2', '3']}
    modification = ArrayDeleteOperation(field='b', value='2')
    return_code = modification.modify(record)
    assert record == {'a': 1, 'b': ['1', '3']}

def test_modification_template():
    record = {'a': '1', 'b': '2'}
    modification = SetOperation(field='c', value='{{ (a | int) + (b | int) }}')
    log.debug("Record before: {}".format(record))
    return_code = modification.modify(record)
    log.debug("Record after: {}".format(record))
    assert record == {'a': '1', 'b': '2', 'c': '3'}

def test_modification_regex_parse():
    record = {'message': 'CRON[12345]: Error during cronjob'}
    modification = RegexParse(field='message', regex='(?P<appname>.*?)\[(?P<pid>\d+)\]: (?P<message>.*)')
    return_code = modification.modify(record)
    assert return_code
    assert record == {'message': 'Error during cronjob', 'appname': 'CRON', 'pid': '12345'}

def test_modification_regex_parse_broken_regex():
    record = {'message': 'CRON[12345]: Error during cronjob'}
    with pytest.raises(ValidationError):
        modification = RegexParse(field='message', regex='(?P<appname.*?)\[(?P<pid>\d+)\]: (?P<message>.*)')

def test_modification_regex_sub():
    record = {'message': 'Error in session 0x2134adf890bc89'}
    modification = RegexSub(field='message', out_field='message', regex='0x[a-fA-F0-9]+', sub='0x###')
    return_code = modification.modify(record)
    assert return_code == True
    assert record == {'message': 'Error in session 0x###'}

class TestKvSet:

    data = {
        'kv': [
            {'dict': 'hostnames-to-owners', 'key': 'ops01', 'value': 'ops'},
            {'dict': 'hostnames-to-owners', 'key': 'ad01', 'value': 'windows'},
            {'dict': 'hostnames-to-owners', 'key': 'linux01', 'value': 'unix'},
        ],
    }

    def test_modify(self, core):
        kv_plugin = core.get_core_plugin('kv')
        record = {'hostname': 'ops01', 'message': 'Issue on ops server'}
        modification = KvSet(dict='hostnames-to-owners', field='hostname', out_field='owner', core=core)
        return_code = modification.modify(record)
        assert return_code == True
        assert record['owner'] == 'ops'
