#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''A module where the snooze core plugin resides'''

from logging import getLogger
from typing import * # pylint: disable=wildcard-import,unused-wildcard-import
from uuid import UUID

import falcon
from pydantic import Field
from falcon import Request, Response

from snooze.plugins.core.basic.plugin import AbortAndWrite, Abort

from snooze.plugins.core import ApiPlugin, ProcessPlugin, BaseApiRoute, parse_uids
from snooze.utils.condition import Condition, AlwaysTrue
from snooze.utils.time_constraints import TemporalConstraint, Anytime
from snooze.utils.model import ApiModel, MongodbMetadata, Record
from snooze.utils.functions import authorize

log = getLogger('snooze.plugins.snooze')

class Snooze(ApiModel):
    '''A snooze filter, which contains conditions to filter incoming alerts'''
    _mongodb = MongodbMetadata(collection='snooze', order_by='name')

    enabled: bool = Field(
        default=True,
        description='When set to false, the snooze filter will be disabled and will not be taken into'
        ' account by the processing',
    )
    name: str = Field(
        default='',
        description='Name of the snooze filter',
    )
    condition: Condition = Field(
        default_factory=AlwaysTrue,
        discriminator='type',
        description='A condition to snooze records matching a given condition',
    )
    time_constraint: TemporalConstraint = Field(
        default_factory=Anytime,
        description='A time constraint that will filter incoming records only if their timestamp match'
        ' the time constraint.',
    )
    comment: str = Field(
        default='',
        description='A comment meant for users, and that describe what this snooze filter is doing',
    )
    discard: bool = Field(
        default=False,
        description='If enabled, will discard all incoming records from the database. '
        'It is meant to be used for alert that create too much logs for snooze to handle/save',
    )
    # Computed
    hits: int = Field(
        default=0,
        description='The number of time a snooze filter matched a record during its lifetime',
    )

    def match(self, record: Record) -> bool:
        '''Whether a record match the Snooze object'''
        return self.condition.match(record) and self.time_constraint.match(record.timestamp)

class SnoozeRetroApplyRoute(BaseApiRoute):
    '''Route used for executing the snooze retro-appy functionality'''
    @authorize
    def on_put(self, _req: Request, resp: Response, uids_str: str):
        '''Retro apply snooze filters, triggered by a PUT /api/snooze_apply/<uid1>+<uid2>+...'''
        uids: List[UUID] = parse_uids(uids_str)
        self.plugin.retro_apply(uids)
        resp.status = falcon.HTTP_OK

class SnoozePlugin(ApiPlugin, ProcessPlugin, model=Snooze):
    '''The snooze process plugin'''
    filters: List[Snooze]

    def __init__(self, *args, **kwargs):
        ApiPlugin.__init__(self, *args, **kwargs)
        ProcessPlugin.__init__(self)
        self.reload()

    def reload(self):
        self.filters = self.database[Snooze].search()

    def load_routes(self, api):
        ApiPlugin.load_routes(self, api)
        api.add_route('/api/snooze_apply/{uids_str}',
            SnoozeRetroApplyRoute(api, self.model, plugin=self))

    def process(self, record: Record) -> Record:
        log.debug("Processing record %s against snooze filters", record.get('hash', ''))
        for snooze in self.filters:
            if snooze.enabled and snooze.match(record):
                log.debug("Snooze %s matched record: %s", snooze.name, record.get('hash', ''))
                record.core.snoozed = snooze.name
                snooze.hits += 1
                self.database[Snooze].increment(dict(uid=snooze.uid), {'hits'})
                self.core.stats.inc('alert_snoozed', {'name': snooze.name})
                if snooze.discard:
                    if record.hash is not None:
                        self.database[Record].delete_by(hash=record.hash)
                    raise Abort()
                raise AbortAndWrite(record)
        return record

    def retro_apply(self, filter_names: List[str]):
        '''Retro applying a list of snooze filters'''
        log.debug("Attempting to retro apply snooze filters %s", filter_names)
        filters = [f for f in self.filters if f.name in filter_names]
        results = []
        for snooze in filters:
            if snooze.enabled and snooze.discard:
                log.debug("Retro apply discard snooze %s", snooze.name)
                result = self.database[Record].delete_search(snooze.condition)
                results.append(result)
