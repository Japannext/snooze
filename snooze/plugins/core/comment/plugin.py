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
from uuid import uuid4, UUID
from enum import Enum
from typing import Any, List, Type
from datetime import datetime

import bson.json_util
import falcon
from falcon import Request, Response
from pydantic import BaseModel, Extra, ValidationError, Field

#from snooze.plugins.core.basic.falcon.route import Route
from snooze.utils.functions import authorize
#from snooze.utils.modification import get_modification
from snooze.utils.typing import Record
from snooze.plugins.core import ApiPlugin, BaseApiRoute, parse_uid, parse_model
from snooze.utils.model import ApiModel, MongodbMetadata
from snooze.utils.typing import SnoozeUser, AuthPayload

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

class Comment(ApiModel):
    '''A data model representing a comment'''
    _mongodb = MongodbMetadata(collection='comment', primaries={'snooze_user'})

    record_uid: UUID = Field(
        required=True,
        description='UUID of the record the comment is associated with',
    )
    date: datetime = Field(
        default_factory=datetime.now,
        description='Date when the comment was written',
    )
    type: CommentType = Field(
        default=CommentType.COMMENT,
        description='Type of comment',
    )
    user: SnoozeUser = Field(
        description='Object representing the user that created the comment',
    )
    message: str = Field(
        default='',
        description='Content of the comment',
    )
    edited: bool = Field(
        default=False,
        description='Indicate if the comment was edited after being created',
    )
    modifications: List[Any] = Field(
        default_factory=list,
        description='A list of modifications to apply to the record. Only used'
        ' when re-escalating, to forcefully re-escalate to another channel, by '
        'overwriting a custom field used by notifications for instance',
    )

class CommentSelfRoute:
    def __init__(self, endpoint):
        self.endpoint = endpoint

    @authorize
    def on_post(self, req: Request, resp: Response):
        '''Create a new comment'''
        document = parse_model(req.media, Comment)
        inject_auth_payload(req.context.auth, req.media)
        self.endpoint.create(document)
        resp.status = falcon.HTTP_CREATED
        resp.media = document.dict()

    @authorize
    def on_put_uid(self, req: Request, resp: Response, uid_str: str):
        uid = parse_uid(uid_str)
        inject_auth_payload(req.context.auth, req.media)
        document = parse_model(req.media, Comment)
        self.endpoint.replace(uid, document)
        resp.status = falcon.HTTP_OK
        resp.media = document.dict()

class PluginMiddleware(Generic[Model], ABC):
    def __init__(self, database, plugin):
        self.database = database
        self.plugin = plugin

    @abstractmethod
    def on_create(self, document: Model):
        ...

    @abstractmethod
    def on_patch(self, search, partial: Partial[Model]):
        ...

    @abstractmethod
    def on_search(self, condition: Condition, pagination: Pagination):
        ...

class CommentMiddleware:
    def on_create(self, comment: Comment):
        self.database[Record].increment(dict(uid=comment.record_uid), {'comment_count'})

class CommentRoute(BaseApiRoute):
    def on_post(self, req: Request, resp: Response):
        BaseApiRoute.on_post(self, req, resp)
        """
        self.database[Record].increment(uid={''})
        record = self.get_associated_record(comment)
        count = len(self.search('comment', ['=', 'record_uid', comment.record_uid])['data'])
        record['comment_count'] = count + 1
        if comment.type.is_active():
            record['state'] = comment.type
        if comment.type.is_modify():
            self.reescalate(comment, record)
        self.update('record', record)
        """

class CommentPlugin(ApiPlugin, model=Comment):
    '''An API plugin to manage comments'''
    notification: Type[ApiPlugin]

    def __init__(self, *args, **kwargs):
        ApiPlugin.__init__(self, *args, **kwargs)
        self.notification = self.core.get_core_plugin('notification')

    def load_routes(self, api):
        ApiPlugin.load_routes(self, api)
        api.add_route('/api/comment_self', CommentSelfRoute(api, Comment, plugin=self))

class OldCommentRoute(Route):
    '''A route for handling comments'''
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.notification_plugin = self.core.get_core_plugin('notification')

    @authorize
    def on_post(self, req, resp):
        for data in req.media:
            try:
                data.setdefault('name', req.context.auth.username)
                comment = Comment(**data)
                log.debug(comment.dict())
            except ValidationError as err:
                raise falcon.HTTPBadRequest(description=f"Invalid comment: {err}") from err
            self.upon_added(comment)
        Route.on_post(self, req, resp)

    @authorize
    def on_put(self, req, resp):
        for data in req.media:
            try:
                data.setdefault('name', req.context.auth.username)
                comment = Comment(**data)
                comment.edited = True
                log.debug(comment.dict())
            except ValidationError as err:
                raise falcon.HTTPBadRequest(description=f"Invalid comment: {err}") from err
            self.upon_edited(comment)
        Route.on_put(self, req, resp)

    @authorize
    def on_delete(self, req, resp, search='[]'):
        if 'uid' in req.params:
            cond_or_uid = ['=', 'uid', req.params['uid']]
        else:
            string = req.params.get('s') or search
            try:
                cond_or_uid = bson.json_util.loads(string)
            except Exception:
                cond_or_uid = string
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
                comments.append(Comment(**data))
            except ValidationError as err:
                log.warning("Invalid comment at uid='%s'", data.get('uid'), exc_info=err)
                continue
        relevant_comments = [x for x in comments if x.type.is_active()]
        if len(relevant_comments) > 0:
            record['state'] = relevant_comments[0].type
        else:
            record['state'] = ''
        self.update('record', record)
