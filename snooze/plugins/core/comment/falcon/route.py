#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Comment custom falcon routes.
Mainly used for allowing users to comment as their user only
'''

from logging import getLogger
from uuid import uuid4
from enum import Enum
from typing import Any, List
from datetime import datetime

import bson.json_util
import falcon
from pydantic import BaseModel, Extra, ValidationError, Field

from snooze.plugins.core.basic.falcon.route import Route
from snooze.utils.functions import authorize
from snooze.utils.modification import get_modification
from snooze.utils.typing import Record

log = getLogger('snooze.api')

class CommentType(str, Enum):
    '''Enum to describe possible comment types'''
    ACK = 'ack'
    ESC = 'esc'
    OPEN = 'open'
    CLOSE = 'close'
    COMMENT = ''

    def is_active(self) -> bool:
        '''Returns if the comment type will modify the record state'''
        return self in [self.ACK, self.ESC, self.OPEN, self.CLOSE]

    def is_modify(self) -> bool:
        '''Returns if the comment type need to modify the record with modifications'''
        return self in [self.ESC, self.OPEN]

class Comment(BaseModel, extra=Extra.allow):
    '''A data model representing a comment'''
    date: str = Field(default_factory=lambda: datetime.now().astimezone().isoformat())
    type: CommentType = CommentType.COMMENT
    record_uid: str
    name: str = 'unknown'
    message: str = ''
    edited: bool = False
    modifications: List[Any] = Field(default_factory=list)

    def __init__(self, data: dict, context=None):
        '''Making sure we get a `name` from the context if not present'''
        if context and 'user' in context and 'name' not in data:
            data['name'] = context['user']['user']['name']
        BaseModel.__init__(self, **data)

class CommentRoute(Route):
    '''A route for handling comments'''
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.notification_plugin = self.core.get_core_plugin('notification')

    @authorize
    def on_post(self, req, resp):
        for data in req.media:
            try:
                comment = Comment(data, req.context)
            except ValidationError as err:
                raise falcon.HTTPBadRequest(description=f"Invalid comment: {err}") from err
            self.upon_added(comment)
        Route.on_post(self, req, resp)

    @authorize
    def on_put(self, req, resp):
        for data in req.media:
            try:
                comment = Comment(data, req.context)
                comment.edited = True
            except ValidationError as err:
                raise falcon.HTTPBadRequest(description=f"Invalid comment: {err}") from err
            self.upon_edited(comment)
        Route.on_put(self, req, resp)

    @authorize
    def on_delete(self, req, resp, search='[]'):
        if 'uid' in req.params:
            cond_or_uid = ['=', 'uid', req.params['uid']]
        else:
            data = req.params.get('s') or search
            cond_or_uid = bson.json_util.loads(data)
        for data in self.search('comment', cond_or_uid)['data']:
            try:
                comment = Comment(**data)
            except ValidationError as err:
                log.warning("Invalid comment at uid='%s'", data.get('uid'), exc_info=err)
                continue
            self.upon_deleted(comment)
        Route.on_delete(self, req, resp, search)

    def reescalate(self, comment: Comment, record: Record):
        '''Re-escalate a record'''
        for mod_data in comment.modifications:
            try:
                modification = get_modification(mod_data, self.core)
                modification.modify(record)
            except Exception as err:
                log.warning("Error in modification %s", modification, exc_info=err)
                continue
        if self.notification_plugin:
            record['notification_plugin'] = {
                'name': comment.name,
                'message': comment.message,
            }
            record.pop('snoozed', None)
            record.pop('notifications', None)
            self.notification_plugin.process(record)

    def get_associated_record(self, comment: Comment) -> dict:
        '''Return the record associated to a comment'''
        records = self.search('record', comment.record_uid)['data']
        if len(records) > 0:
            return records[0]
        else:
            raise falcon.HTTPInternalServerError(
                description=f"No record found for comment[record_uid] = '{comment.record_uid}'")

    def upon_added(self, comment: Comment):
        '''Update the record when one of its comment is added'''
        record = self.get_associated_record(comment)
        count = len(self.search('comment', ['=', 'record_uid', comment.record_uid])['data'])
        record['comment_count'] = count + 1
        if comment.type.is_active():
            record['state'] = comment.type
        if comment.type.is_modify():
            self.reescalate(comment, record)
        self.update('record', record)

    def upon_edited(self, comment: Comment):
        '''Update the record when one of its comment is edited'''
        record = self.get_associated_record(comment)
        if comment.type.is_active():
            record['state'] = comment.type
        if comment.type.is_modify():
            self.reescalate(comment, record)
        self.update('record', record)

    def upon_deleted(self, comment: Comment):
        '''Update the record when one of its comment is deleted'''
        record = self.get_associated_record(comment)
        count = len(self.search('comment', ['=', 'record_uid', comment.record_uid])['data'])
        record['comment_count'] = count - 1
        comments = []
        for data in self.search('comment', ['=', 'record_uid', comment.record_uid])['data']:
            try:
                comments.append(Comment(data))
            except ValidationError as err:
                log.warning("Invalid comment at uid='%s'", data.get('uid'), exc_info=err)
                continue
        relevant_comments = [x for x in comments if x.is_active()]
        if len(relevant_comments) > 0:
            record['state'] = relevant_comments[0].type
        else:
            record['state'] = ''
        self.update('record', record)
