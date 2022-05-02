#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

from snooze.utils.parser import parser

class TestParserLogic:
    def test_word(self):
        result = parser('hello')
        assert result.dict() == {'type': 'SEARCH', 'value': 'hello'}

    def test_key_value(self):
        result = parser('key = value')
        assert result.dict() == {'type': '=', 'field': 'key', 'value': 'value'}

    def test_and(self):
        result = parser('key1=value1 AND key2=value2')
        assert result.dict() == {
            'type': 'AND',
            'conditions': [
                {'type': '=', 'field': 'key1', 'value': 'value1'},
                {'type': '=', 'field': 'key2', 'value': 'value2'},
            ],
        }

    def test_and_symbol(self):
        result = parser('key1=value1&key2=value2')
        assert result.dict() == {
            'type': 'AND',
            'conditions': [
                {'type': '=', 'field': 'key1', 'value': 'value1'},
                {'type': '=', 'field': 'key2', 'value': 'value2'},
            ],
        }

    def test_implicit_and(self):
        result = parser('key1=value1 key2=value2')
        assert result.dict() == {
            'type': 'AND',
            'conditions': [
                {'type': '=', 'field': 'key1', 'value': 'value1'},
                {'type': '=', 'field': 'key2', 'value': 'value2'},
            ],
        }

    def test_or(self):
        result = parser('key1=value1 OR key2=value2')
        assert result.dict() == {
            'type': 'OR',
            'conditions': [
                {'type': '=', 'field': 'key1', 'value': 'value1'},
                {'type': '=', 'field': 'key2', 'value': 'value2'},
            ],
        }

    def test_or_symbol(self):
        result = parser('key1=value1|key2=value2')
        assert result.dict() == {
            'type': 'OR',
            'conditions': [
                {'type': '=', 'field': 'key1', 'value': 'value1'},
                {'type': '=', 'field': 'key2', 'value': 'value2'},
            ],
        }

    def test_not(self):
        result = parser('not key1=value1')
        assert result.dict() == {'type': 'NOT', 'condition': {'type': '=', 'field': 'key1', 'value': 'value1'}}

    def test_not_symbol(self):
        result = parser('!key1=value1')
        assert result.dict() == {'type': 'NOT', 'condition': {'type': '=', 'field': 'key1', 'value': 'value1'}}

    def test_double_negation(self):
        result = parser('!!key1=value1')
        assert result.dict() == {'type': 'NOT', 'condition': {'type': 'NOT', 'condition': {'type': '=', 'field': 'key1', 'value': 'value1'}}}

    def test_parenthesis(self):
        result = parser('NOT (key1=value1 AND key2=value2)')
        assert result.dict() == {
            'type': 'NOT',
            'condition': {
                'type': 'AND',
                'conditions': [
                    {'type': '=', 'field': 'key1', 'value': 'value1'},
                    {'type': '=', 'field': 'key2', 'value': 'value2'},
                ],
            }
        }

    def test_priority(self):
        result = parser('NOT key1=value1 AND key2=value2')
        assert result.dict() == {
            'type': 'AND',
            'conditions': [
                {'type': 'NOT', 'condition': {'type': '=', 'field': 'key1', 'value': 'value1'}},
                {'type': '=', 'field': 'key2', 'value': 'value2'},
            ],
        }

    def test_complex_query(self):
        result = parser('myapp and source=syslog and custom_field = "myapp01" or custom_field = "myapp02"')
        assert result.dict() == {
            'type': 'AND',
            'conditions': [
                {'type': 'SEARCH', 'value': 'myapp'},
                {'type': '=', 'field': 'source', 'value': 'syslog'},
                {
                    'type': 'OR',
                    'conditions': [
                        {'type': '=', 'field': 'custom_field', 'value': 'myapp01'},
                        {'type': '=', 'field': 'custom_field', 'value': 'myapp02'},
                    ],
                }
            ],
        }

    def test_complex_query_parenthesis(self):
        result = parser('myapp and source=syslog and (custom_field = "myapp01" or custom_field = "myapp02")')
        assert result.dict() == {
            'type': 'AND',
            'conditions': [
                {'type': 'SEARCH', 'value': 'myapp'},
                {'type': '=', 'field': 'source', 'value': 'syslog'},
                {
                    'type': 'OR',
                    'conditions': [
                        {'type': '=', 'field': 'custom_field', 'value': 'myapp01'},
                        {'type': '=', 'field': 'custom_field', 'value': 'myapp02'},
                    ],
                }
            ],
        }

class TestParserTypes:
    def test_integer(self):
        result = parser('a = 123')
        assert result.dict() == {'type': '=', 'field': 'a', 'value': 123}

    def test_negative_integer(self):
        result = parser('a = -42')
        assert result.dict() == {'type': '=', 'field': 'a', 'value': -42}

    def test_float(self):
        result = parser('a = 3.14')
        assert result.dict() == {'type': '=', 'field': 'a', 'value': 3.14}

    def test_bool_true(self):
        result = parser('mybool=true')
        assert result.dict() == {'type': '=', 'field': 'mybool', 'value': True}

    def test_bool_false(self):
        result = parser('mybool=false')
        assert result.dict() == {'type': '=', 'field': 'mybool', 'value': False}

    def test_double_quoted_string(self):
        result = parser('key = "value"')
        assert result.dict() == {'type': '=', 'field': 'key', 'value': 'value'}

    def test_single_quoted_string(self):
        result = parser("key = 'value'")
        assert result.dict() == {'type': '=', 'field': 'key', 'value': 'value'}

    def test_double_quoted_escape(self):
        result = parser(r'key = "value\t\n\\"')
        assert result.dict() == {'type': '=', 'field': 'key', 'value': "value\t\n\\"}

    def test_double_quoted_escape_quote(self):
        result = parser(r'key = "my \"test\""')
        assert result.dict() == {'type': '=', 'field': 'key', 'value': 'my "test"'}

    def test_quoted_field(self):
        result = parser('"myfield with space" = myapp01')
        assert result.dict() == {'type': '=', 'field': 'myfield with space', 'value': 'myapp01'}

    def test_single_quote_string(self):
        result = parser("'myfield' = 'myvalue'")
        assert result.dict() == {'type': '=', 'field': 'myfield', 'value': 'myvalue'}

    def test_array(self):
        result = parser('myfield = [1, 2, 3]')
        assert result.dict() == {'type': '=', 'field': 'myfield', 'value': [1, 2, 3]}

    def test_nested_array(self):
        result = parser('myfield = [[1], [2], [3]]')
        assert result.dict() == {'type': '=', 'field': 'myfield', 'value': [[1], [2], [3]]}

    def test_dict(self):
        result = parser('myfield = {a: 1, b: 2}')
        assert result.dict() == {'type': '=', 'field': 'myfield', 'value': {'a': 1, 'b': 2}}

    def test_nested_dict(self):
        result = parser('myfield = {a: {"mymessage": "x"}, b: 2}')
        assert result.dict() == {'type': '=', 'field': 'myfield', 'value': {'a': {'mymessage': 'x'}, 'b': 2}}

    def test_hash(self):
        result = parser('hash=3f75728488a0e6892905f0db6a473382')
        assert result.dict() == {'type': '=', 'field': 'hash', 'value': '3f75728488a0e6892905f0db6a473382'}

class TestParserOperations:
    def test_nequal(self):
        result = parser('process != systemd')
        assert result.dict() == {'type': '!=', 'field': 'process', 'value': 'systemd'}

    def test_matches(self):
        result = parser('message MATCHES "[aA]lert"')
        assert result.dict() == {'type': 'MATCHES', 'field': 'message', 'value': '[aA]lert'}

    def test_matches_symbol(self):
        result = parser('message ~ "[aA]lert"')
        assert result.dict() == {'type': 'MATCHES', 'field': 'message', 'value': '[aA]lert'}

    def test_exists(self):
        result = parser('custom_field EXISTS')
        assert result.dict() == {'type': 'EXISTS', 'field': 'custom_field'}

    def test_exists_symbol(self):
        result = parser('custom_field?')
        assert result.dict() == {'type': 'EXISTS', 'field': 'custom_field'}

    def test_exists_symbol_not(self):
        result = parser('!custom_field?')
        assert result.dict() == {'type': 'NOT', 'condition': {'type': 'EXISTS', 'field': 'custom_field'}}

    def test_gt(self):
        result = parser('mail_queue>100')
        assert result.dict() == {'type': '>', 'field': 'mail_queue', 'value': 100}

    def test_lt(self):
        result = parser('port < 1024')
        assert result.dict() == {'type': '<', 'field': 'port', 'value': 1024}

    def test_contains(self):
        result = parser('rules contains myrule')
        assert result.dict() == {'type': 'CONTAINS', 'field': 'rules', 'value': 'myrule'}

    def test_contains_array(self):
        result = parser('myarray contains [1, 2, 3]')
        assert result.dict() == {'type': 'CONTAINS', 'field': 'myarray', 'value': [1, 2, 3]}

    def test_in(self):
        result = parser('rule in [rule01, rule02, rule03]')
        assert result.dict() == {'type': 'IN', 'field': 'rule', 'value': ['rule01', 'rule02', 'rule03']}
