#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''A plugin for defining rules to apply modifications on records matching the
rule's condition'''

from __future__ import annotations

from logging import getLogger
from typing import List, Dict, Optional
from uuid import UUID

from pydantic import Field

from snooze.database import Pagination
from snooze.plugins.core import ApiPlugin, ProcessPlugin
from snooze.utils.condition import Condition, AlwaysTrue, InArray, Exists
from snooze.utils.model import ApiModel, MongodbMetadata, Record
from snooze.utils.modification import Modification

log = getLogger('snooze.process')

class Rule(ApiModel):
    '''Apply modification to records based on conditions'''
    _mongodb = MongodbMetadata(collection='rule')

    enabled: bool = Field(
        default=True,
        description='If set to false, disable the rule and ignore it during record processing',
    )
    name: str = Field(
        default='',
        description='The name of the rule',
    )
    condition: Condition = Field(
        discriminator='type',
        default_factory=AlwaysTrue,
        description='Apply the rule modifications and childrens to all records matching the condition',
    )
    modifications: List[Modification] = Field(
        default_factory=list,
        description='A list of modifications to apply to the record matching the rule',
    )
    comment: str = Field(
        default='',
        description='Description of what the rule is doing',
    )
    parent: Optional[UUID] = Field(
        default=None,
        description='UUID of the parent if the rule is nested. If absent, this will be a top-level rule',
    )

    def match(self, record: Record) -> bool:
        '''Check if a record matched this rule's condition'''
        match = self.condition.match(record)
        if match:
            record.core.rules.append(self.name)
        return match

    def modify(self, record: Record) -> bool:
        '''Modify the record based of this rule's modifications'''
        modified = False
        modifs = []
        for modification in self.modifications:
            if modification.modify(record):
                modified = True
                modifs.append(modification)
        if modified:
            log.debug("Record %s has been modified: %s", record.get('hash', ''), [m.pprint() for m in modifs])
        else:
            log.debug("Record %s has not been modified", record.get('hash', ''))
        return modified

class RulePlugin(ApiPlugin, ProcessPlugin, model=Rule):
    '''The rule plugin main class'''
    rules: List[Rule]
    childrens: Dict[UUID, List[Rule]]

    def __init__(self, *args, **kwargs):
        ApiPlugin.__init__(self, *args, **kwargs)
        ProcessPlugin.__init__(self)
        self.reload()

    def reload(self):
        log.debug("Reloading data for plugin %s", self.name)
        self.rules = self.database[Rule].search(~Exists(field='parent'), Pagination(order_by='name'))
        self.build_childrens(self.rules)

    def build_childrens(self, uids: List[UUID]):
        '''Build a dictionary mapping rule uids to their childrens'''
        children_uids: List[UUID] = []
        for children in self.database[Rule].search(InArray(field='parent', value=uids)):
            uid = children.parent_uid
            self.childrens.setdefault(uid, [])
            self.childrens[uid].append(children)
            children_uids.append(children.uid)
        self.build_childrens(children_uids)

    def process(self, record):
        '''Process the record against a list of rules'''
        self.process_rules(record, self.rules)
        return record

    def process_rules(self, record, rules):
        '''Process a list of rules'''
        log.debug("Processing record %s against rules", record.get('hash', ''))
        for rule in rules:
            if rule.enabled and rule.match(record):
                log.debug("Rule %s matched record: %s", rule.name, record.get('hash', ''))
                rule.modify(record)
                self.process_rules(record, self.childrens.get(rule.uid, []))
