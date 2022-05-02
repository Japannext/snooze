#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#
'''
Tests for conditions. Each condition is tested in its own class.
'''

import pytest
from pydantic import ValidationError

from snooze.utils.condition import *

class TestEquals:
    def test_guess(self):
        condition = guess_condition({'type': '=', 'field': 'a', 'value': '0'})
        assert isinstance(condition, Equals)
        assert condition.type == '='
    def test_init(self):
        condition = Equals(field='a', value='0')
        assert isinstance(condition, ConditionBase)
    def test_match_simple(self):
        record = {'a': '1', 'b': '2'}
        condition = Equals(field='a', value='1')
        assert condition.match(record) == True
    def test_match_nested_dict(self):
        record = {'a': '1', 'b': {'c': '1'}}
        condition = Equals(field='b.c', value='1')
        assert condition.match(record) == True
    def test_miss_nested_dict(self):
        record = {'a': '1', 'b': {'c': 1}}
        condition = Equals(field='a.c', value='2')
        assert condition.match(record) == False
    def test_match_nested_list(self):
        record = {'a': ['1', '2']}
        condition = Equals(field='a.1', value='2')
        assert condition.match(record) == True
    def test_str(self):
        condition = Equals(field='a.1', value='2')
        assert str(condition) == "a.1 = '2'"
    def test_edge_no_field(self):
        record = {'a': '1'}
        with pytest.raises(ValidationError):
            condition = Equals(field=None, value='1')
    def test_edge_no_value(self):
        record = {'a': '1'}
        condition = Equals(field='a', value=None)
        assert condition.match(record) == False

class TestNotEquals:
    def test_guess_condition(self):
        condition = guess_condition({'type': '!=', 'field': 'a', 'value': '0'})
        assert isinstance(condition, NotEquals)
        assert condition.type == '!='
    def test_init(self):
        condition = NotEquals(field='a', value='1')
        assert isinstance(condition, ConditionBase)
    def test_miss(self):
        record = {'a': '1', 'b': '2'}
        condition = NotEquals(field='a', value='1')
        assert condition.match(record) == False
    def test_str(self):
        condition = NotEquals(field='a', value=1)
        assert str(condition) == "a != 1"

class TestGreaterThan:
    def test_guess_condition(self):
        condition = guess_condition({'type': '>', 'field': 'a', 'value': '0'})
        assert isinstance(condition, GreaterThan)
    def test_init(self):
        condition = GreaterThan(field='b', value='1')
        assert isinstance(condition, ConditionBase)
    def test_match_two_float(self):
        record = {'a': 1.0, 'b': 2.0}
        condition = GreaterThan(field='b', value=1.0)
        assert condition.match(record) == True
    def test_match_string_and_integer(self):
        record = {'a': 1, 'b': 2}
        condition = GreaterThan(field='b', value='1')
        assert condition.match(record) == False
    def test_str(self):
        condition = GreaterThan(field='x', value=100)
        assert str(condition) == "x > 100"

class TestLowerThan:
    def test_guess_condition(self):
        condition = guess_condition({'type': '<', 'field': 'a', 'value': '0'})
        assert isinstance(condition, LowerThan)
    def test_init(self):
        condition = LowerThan(field='var', value='ab')
        assert isinstance(condition, ConditionBase)
    def test_match_two_string(self):
        record = {'var': 'aa'}
        condition = LowerThan(field='var', value='ab')
        assert condition.match(record) == True
    def test_str(self):
        condition = LowerThan(field='x', value=100)
        assert str(condition) == "x < 100"

class TestAnd:
    def test_guess_condition(self):
        condition = guess_condition({'type': 'AND', 'conditions': [
            {'type': '=', 'field': 'a', 'value': 1},
            {'type': '=', 'field': 'b', 'value': 2},
        ]})
        assert isinstance(condition, And)
    def test_init(self):
        condition = And(conditions=[Equals(field='a', value=1), Equals(field='b', value=2)])
        assert isinstance(condition, ConditionBase)
    def test_operation(self):
        condition = Equals(field='a', value=1) & Equals(field='b', value=2)
        assert isinstance(condition, And)
    def test_matches(self):
        record = {'a': 1, 'b': 2}
        conditions = [
            Equals(field='a', value=1) & Equals(field='b', value=2)
        ]
        for condition in conditions:
            assert condition.match(record) == True
    def test_misses(self):
        record = {'a': 1, 'b': 2}
        conditions = [
            Equals(field='a', value=1) & Equals(field='b', value=3)
        ]
        for condition in conditions:
            assert condition.match(record) is False
    def test_multiple(self):
        record = {'a': 1, 'b': 2, 'c': 3}
        condition = And(conditions=[
            Equals(field='a', value=1),
            Equals(field='b', value=2),
            Equals(field='c', value=3),
        ])
        assert condition.match(record) is True
    def test_nested(self):
        record = {'a': 1, 'b': 2, 'c': 3}
        condition = Equals(field='a', value=1) & (Equals(field='b', value=2) & Equals(field='c', value=3))
        assert condition.match(record) is True
    def test_nested_miss(self):
        record = {'a': 1, 'b': 2, 'c': 3}
        condition = Equals(field='a', value=1) & (Equals(field='b', value=2) & Equals(field='c', value=4))
        assert condition.match(record) is False
    def test_str(self):
        condition = Equals(field='a', value=1) & Equals(field='b', value=2)
        assert str(condition) == "(a = 1 & b = 2)"

class TestOr:
    def test_guess_condition(self):
        condition = guess_condition({'type': 'OR', 'conditions': [Equals(field='a', value=1), Equals(field='b', value=2)]})
        assert isinstance(condition, Or)
    def test_init(self):
        condition = Or(conditions=[Equals(field='a', value=1), Equals(field='b', value=2)])
        assert isinstance(condition, ConditionBase)
    def test_match(self):
        record = {'a': 1, 'b': 3}
        condition = Or(conditions=[Equals(field='a', value=1), Equals(field='b', value=2)])
        assert condition.match(record) is True
    def test_multiple(self):
        record = {'a': 1, 'b': 2, 'c': 3}
        condition = Or(conditions=[Equals(field='a', value=6), Equals(field='b', value=4), Equals(field='c', value=3)])
        assert condition.match(record) is True
    def test_str(self):
        condition = Equals(field='a', value=1) | Equals(field='b', value=2)
        assert str(condition) == "(a = 1 | b = 2)"

class TestNot:
    def test_guess_condition(self):
        condition = guess_condition({'type': 'NOT', 'condition': Equals(field='a', value=1)})
        assert isinstance(condition, Not)
    def test_init(self):
        condition = Not(condition=Equals(field='a', value=1))
        assert isinstance(condition, ConditionBase)
    def test_match(self):
        record = {'a': 1, 'b': 3}
        condition = Not(condition=Equals(field='a', value=2))
        assert condition.match(record) is True
    def test_miss(self):
        record = {'a': 1, 'b': 3}
        condition = Not(condition=Equals(field='a', value=1))
        assert condition.match(record) is False
    def test_str(self):
        condition = ~Equals(field='b', value=2)
        assert str(condition) == "!(b = 2)"

class TestMatches:
    def test_guess_condition(self):
        condition = guess_condition({'type': 'MATCHES', 'field': 'a', 'value': 'string'})
        assert isinstance(condition, Matches)
    def test_match(self):
        record = {'a': '__pattern__'}
        condition = Matches(field='a', value='pattern')
        assert condition.match(record) is True
    def test_match_sugar(self):
        record = {'a': '__pattern__'}
        condition = Matches(field='a', value='/pattern/')
        assert condition.match(record) is True
    def test_str(self):
        condition = Matches(field='a', value='/string/')
        assert str(condition) == "a ~ /string/"

class TestExists:
    def test_guess_condition(self):
        condition = guess_condition({'type': 'EXISTS', 'field': 'a'})
        assert isinstance(condition, Exists)
    def test_match(self):
        record = {'a': '1'}
        condition = Exists(field='a')
        assert condition.match(record) == True
    def test_miss(self):
        record = {'a': '1'}
        condition = Exists(field='b')
        assert condition.match(record) == False
    def test_str(self):
        condition = Exists(field='a')
        assert str(condition) == "a?"

class TestContains:
    def test_guess_condition(self):
        condition = guess_condition({'type': 'CONTAINS', 'field': 'a', 'value': 'element'})
        assert isinstance(condition, ArrayContains)
    """
    def test_match_search_in_string(self):
        record = {'a': ['0', ['11', '2', 9], '3']}
        condition1 = Contains(['CONTAINS', 'a', '1'])
        condition2 = Contains(['CONTAINS', 'a', 9])
        assert condition1.match(record) == True
        assert condition2.match(record) == True
    def test_match_incomplete_list(self):
        record = {'a': '11', 'b': 9}
        condition1 = Contains(['CONTAINS', 'a', ['0', '1']])
        condition2 = Contains(['CONTAINS', 'b', ['0', 9]])
        assert condition1.match(record) == True
        assert condition2.match(record) == True
    """
    def test_str(self):
        condition = ArrayContains(field='a', value='element')
        assert str(condition) == "a contains 'element'"

"""
class TestIn:
    def test_guess_condition(self):
        condition = guess_condition(['IN', ['1', '2', '3'], 'a'])
        assert isinstance(condition, In)
    def test_match_list(self):
        record = {'a': '1', 'b': 1}
        condition1 = In(['IN', ['1', '5'], 'a'])
        condition2 = In(['IN', [1, 5], 'b'])
        assert condition1.match(record) == True
        assert condition2.match(record) == True
    def test_miss_list(self):
        record = {'a': ['0', ['11', '2'], '3']}
        condition = In(['IN', ['1', '5'], 'a'])
        assert condition.match(record) == False
    def test_match_condition(self):
        record = {'a': [{'b':'0'}, {'c': '0'}]}
        condition = In(['IN', ['=', 'c', '0'], 'a'])
        assert condition.match(record) == True
    def test_miss_condition(self):
        record = {'a': [{'b':'0'}, {'c': '0'}]}
        condition = In(['IN', ['=', 'd', '0'], 'a'])
        assert condition.match(record) == False
    def test_match_integer(self):
        record = {'a': [{'b':0}, {'c': '0'}]}
        condition = In(['IN', ['=', 'b', 0], 'a'])
        assert condition.match(record) == True
    def test_str(self):
        condition = In(['IN', 'element', 'a'])
        assert str(condition) == "('element' in a)"
"""

class TestSearch:
    def test_guess_condition(self):
        condition = guess_condition({'type': 'SEARCH', 'value': 'string'})
        assert isinstance(condition, Search)
    def test_match_incomplete_field(self):
        record = {'myfield': [{'b':'mystring'}, {'mysearch': '0'}]}
        condition = Search(value='field')
        assert condition.match(record) is True
    def test_match_nested_value(self):
        record = {'myfield': [{'b':'mystring'}, {'mysearch': '0'}]}
        condition = Search(value='string')
        assert condition.match(record) is True
    def test_match_incomplete_nested_field(self):
        record = {'myfield': [{'b':'mystring'}, {'mysearch': '0'}]}
        condition = Search(value='search')
        assert condition.match(record) is True
    def test_miss(self):
        record = {'myfield': [{'b':'mystring'}, {'mysearch': '0'}]}
        condition = Search(value='value')
        assert condition.match(record) is False
    def test_str(self):
        condition = Search(value='mystring')
        assert str(condition) == "(SEARCH 'mystring')"

class TestAlwaysTrue:
    def test_guess_condition(self):
        condition = guess_condition({'type': 'ALWAYS_TRUE'})
        assert isinstance(condition, AlwaysTrue)
        assert condition.type == 'ALWAYS_TRUE'

"""
class TestAlwaysTrue:
    def test_guess_condition(self):
        conditions = [
            [],
            [''],
            [None],
            None,
        ]
        for condition_str in conditions:
            condition = guess_condition(condition_str)
            assert isinstance(condition, AlwaysTrue)

    def test_match(self):
        records = [
            {'a': 1},
            {'b': '2'},
            {},
        ]
        condition = AlwaysTrue()
        for record in records:
            assert condition.match(record) == True
"""
