#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Combinatory parser for the query language'''

import pyparsing as pp
from pyparsing import pyparsing_common as ppc

from pydantic import BaseModel, Field

from snooze.utils.condition import Condition, And, Or, Not, guess_condition

pp.ParserElement.enablePackrat()

EQUAL = pp.Literal('=')
NOT_EQUAL = pp.Literal('!=')
GT = pp.Literal('>')
LT = pp.Literal('<')
GTE = pp.Literal('>=')
LTE = pp.Literal('<=')
MATCHES = (pp.CaselessKeyword('MATCHES') | pp.Literal('~')).setParseAction(lambda: 'MATCHES')
EXISTS = (pp.CaselessKeyword('EXISTS') | pp.Literal('?')).setParseAction(lambda: 'EXISTS')
CONTAINS = (pp.CaselessKeyword('CONTAINS')).setParseAction(lambda: 'CONTAINS')
IN = (pp.CaselessKeyword('IN')).setParseAction(lambda: 'IN')
SEARCH = pp.CaselessKeyword('SEARCH')

LPAR, RPAR, LBRACK, RBRACK, LBRACE, RBRACE, COLON = map(pp.Suppress, '()[]{}:')

NOT = (pp.CaselessKeyword('NOT') | '!').setParseAction(lambda: 'NOT')
OR = (pp.CaselessKeyword('OR') | '|').setParseAction(lambda: 'OR')
AND = pp.Optional(pp.CaselessKeyword('AND') | '&').setParseAction(lambda: 'AND')

valid_word = pp.Regex(r'[a-zA-Z0-9_.-]+')

string = pp.QuotedString('"', escChar="\\") | pp.QuotedString("'", escChar="\\")
boolean = (
    pp.CaselessKeyword('true').setParseAction(lambda: True)
    | pp.CaselessKeyword('false').setParseAction(lambda: False)
)

literal = pp.Forward()

array_elements = pp.delimitedList(literal, delim=',')
array = pp.Group(LBRACK + pp.Optional(array_elements, []) + RBRACK)
array.setParseAction(lambda t: t.asList())
hashmap = pp.Forward()

fieldname = string | valid_word
literal << (ppc.real ^ ppc.signed_integer ^ string ^ array ^ hashmap ^ boolean ^ valid_word)

hashmap_element = pp.Group(fieldname + COLON + literal)
hashmap_elements = pp.delimitedList(hashmap_element, delim=',')
hashmap << pp.Dict(LBRACE + pp.Optional(hashmap_elements) + RBRACE)
hashmap.setParseAction(lambda t: t.asDict())

term = pp.Forward()
expression = pp.Forward()

binary_operation = EQUAL | NOT_EQUAL | MATCHES | GT | LT | GTE | LTE | CONTAINS | IN
field_operation = EXISTS
unary_operation = SEARCH

term << (
    (fieldname('field') + binary_operation('type') + literal('value'))
    | (fieldname('field') + field_operation('type'))
    | (unary_operation('type') + fieldname('value'))
    | literal('value')
)

def term_parser(tokens):
    '''Parse an expression'''
    if 'type' not in tokens:
        tokens['type'] = 'SEARCH'
    return guess_condition(tokens)

term.setParseAction(term_parser)

# Parse expressions that have an order of priority in operations
logic_operations = [
    (NOT, 1, pp.opAssoc.RIGHT, lambda tokens: Not(condition=tokens[0][1])),
    (OR, 2, pp.opAssoc.LEFT, lambda tokens: Or(conditions=tokens[0][::2])),
    (AND, 2, pp.opAssoc.LEFT, lambda tokens: And(conditions=tokens[0][::2])),
]

expression << pp.infixNotation(term, logic_operations)

def parser(data: str) -> Condition:
    '''Parse a query string'''
    result = expression.parseString(data)[0]
    return result
