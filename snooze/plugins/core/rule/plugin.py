#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''A plugin for defining rules to apply modifications on records matching the
rule's condition'''

from logging import getLogger
from typing import List

from snooze.plugins.core import Plugin
from snooze.utils.condition import get_condition, validate_condition
from snooze.utils.modification import get_modification, validate_modification
from snooze.utils.typing import Record, Rule as RuleType

apilog = getLogger('snooze.api')
proclog = getLogger('snooze.process')

class Rule(Plugin):
    '''The rule plugin main class'''
    def process(self, record):
        '''Process the record against a list of rules'''
        context = dict(rid=record.get('uid', 'none'), plugin=self.name)
        proclog.debug('Start', extra=context)
        self.process_rules(record, self.rules)
        proclog.debug('Done', extra=context)
        return record

    def validate(self, obj):
        '''Validate a rule object'''
        validate_condition(obj)
        validate_modification(obj, self.core)

    def process_rules(self, record, rules):
        '''Process a list of rules'''
        context = dict(rid=record.get('uid', 'none'), plugin=self.name)
        for rule in rules:
            if rule.enabled and rule.match(record):
                proclog.debug("Rule '%s' (%s) matched record", rule.name, rule.uid, extra=context)
                rule.modify(record)
                self.process_rules(record, rule.children)

    def reload_data(self, sync = False):
        context = dict(plugin=self.name)
        apilog.info("Reloading...", extra=context)
        self.data = self.db.search('rule', ['NOT', ['EXISTS', 'parent']], orderby=self.meta.force_order)['data']
        rules = []
        for rule in (self.data or []):
            rules.append(RuleObject(rule, self))
        self.rules = rules
        if sync:
            self.sync_neighbors()
        apilog.info("Reload successful", extra=context)

class RuleObject:
    '''An object representing the rule object in the database'''
    def __init__(self, rule: RuleType, rule_plugin: Rule = None):
        uid = rule.get('uid')
        context = dict(plugin='rule', uid=uid)
        core = None
        order = None
        if rule_plugin:
            core = rule_plugin.core
            order = rule_plugin.meta.force_order
        self.enabled = rule.get('enabled', True)
        self.name = rule['name']
        self.uid = uid
        apilog.debug("Instantiating rule: %s", self.name, extra=context)
        self.condition = get_condition(rule.get('condition'))
        apilog.debug("Instantiating condition: %s", self.condition, extra=context)
        self.modifications = []
        for modification in (rule.get('modifications') or []):
            apilog.debug("Instantiating modification: %s", modification, extra=context)
            self.modifications.append(get_modification(modification, core=core))
        apilog.debug('Searching children', extra=context)
        self.children = []
        if core and core.db:
            db = core.db
            children = db.search('rule', ['=', 'parent', uid], orderby=order)['data']
            for child_rule in children:
                apilog.debug("Found child '%s'", child_rule['name'], extra=context)
                self.children.append(RuleObject(child_rule, rule_plugin))

    def match(self, record: Record) -> bool:
        '''Check if a record matched this rule's condition'''
        match = self.condition.match(record)
        if match:
            if not 'rules' in record:
                record['rules'] = []
            if self.name not in record['rules']:
                record['rules'].append(self.name)
        return match

    def modify(self, record: Record) -> bool:
        '''Modify the record based of this rule's modifications'''
        context = dict(plugin='rule', rid=record.get('uid'))
        modified = False
        for modification in self.modifications:
            if modification.modify(record):
                modified = True
                proclog.info("Record modified by rule '%s': %s", self.name, modification.pprint(), extra=context)
        if not modified:
            proclog.debug("Record not modified by rule '%s'", self.name, extra=context)
        return modified

    def __repr__(self):
        return self.name
