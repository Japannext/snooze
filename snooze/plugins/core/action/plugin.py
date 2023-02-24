#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''The action plugin. Execute an async action (script or python code) upon being
triggered by a notification. Actions can be batched or delayed.'''

import socket
import time
from threading import Event
from typing import Optional

from logging import getLogger

from snooze.plugins.core import Plugin
from snooze.utils.functions import ensure_hash
from snooze.utils.mq import Worker
from snooze.utils.threading import SurvivingThread

proclog = getLogger('snooze-process')
apilog = getLogger('snooze-api')

class Action(Plugin):
    '''The action plugin. Spawn a background thread that will manage delayed and batched actions'''
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.hostname = socket.gethostname()
        self.actions = []
        self.delayed_actions = DelayedActions(self, self.core.exit_event)
        self.core.threads['delayed_actions'] = self.delayed_actions

    def post_init(self):
        self.reload_data()

        delayed_actions = self.core.db.search('action.delay', ['=', 'host', self.hostname])
        if delayed_actions['count'] > 0:
            for action_obj in delayed_actions['data']:
                action_uid = action_obj['action_uid']
                queue_it = [action for action in self.actions if action.uid == action_uid]
                if len(queue_it) > 0:
                    action = queue_it[0]
                    action_obj['action'] = action
                    action_obj['time'] = time.time() + action_obj['delay']
                    action_obj['record'] = {'hash': action_obj['record_hash']}
                    self.delayed_actions.set_delayed(action_obj, False)
                else:
                    apilog.warning("Delayed notification '%s' original notification in not in the database anymore. "
                        "Removing it from queue", action_obj)
                    self.core.db.delete('action.delay', ['=', 'uid', action_uid])
            apilog.debug("Restored delayed actions %s", self.delayed_actions.delayed)

    def _post_reload(self):
        actions = []
        for action in self.data or []:
            action_object = ActionObject(action, self)
            actions.append(action_object)
            if action_object.batch:
                self.core.mq.update_queue(f"action_{action_object.uid}",
                    action_object.batch_timer, action_object.batch_maxsize, ActionWorker, action_object)
            apilog.debug("Init action %s", action_object.name, extra=dict(plugin=self.name))
        self.core.mq.keep_queues([f"action_{action.uid}" for action in actions], "action_")
        self.actions = actions
        notification_plugin = self.core.get_core_plugin('notification')
        if notification_plugin:
            notification_plugin.reload_data()

class UnknownPlugin(RuntimeError):
    '''A runtime error raised when the action plugin is not found. Useful when
    the action plugin comes from an external pip package'''
    def __init__(self, plugin_name: str):
        RuntimeError.__init__(self, f"Plugin '{plugin_name}' not found")

class ActionObject:
    '''Object representing an action in the database'''
    def __init__(self, action, plugin):
        self.action = action
        self.plugin = plugin
        self.delayed = plugin.delayed_actions
        self.core = plugin.core
        self.uid = action.get('uid')
        self.name = action.get('name', '')
        self.selected = action.get('action', {}).get('selected', '')
        self.content = action.get('action', {}).get('subcontent', {})
        self.content['action_name'] = self.name
        self.action_plugin = self.core.get_core_plugin(self.selected)
        if self.action_plugin is None:
            raise UnknownPlugin(self.selected)
        batch = self.action_plugin.meta.options.get('batch', {})
        self.batch = self.content.get('batch', batch.get('default', False))
        self.content['batch'] = self.batch
        self.batch_timer = self.content.get('batch_timer', batch.get('timer', 10))
        self.batch_maxsize = self.content.get('batch_maxsize', batch.get('maxsize', 100))

    def send(self, action_obj):
        '''Called when the action is triggered'''
        record = action_obj['record']
        if action_obj['delay'] <= 0 and action_obj['total'] != 0:
            if action_obj['every'] <= 0:
                if action_obj['total'] < 0:
                    proclog.warning("Action '%s' is probably misonfigured (spamming). Will send only once.",
                        self.name)
                    loop = 1
                else:
                    loop = action_obj['total']
            else:
                loop = 1
            proclog.info("Executing action '%s' (%d times)", self.name, loop)
            self.send_one(loop, action_obj)
        if action_obj['total'] != 0:
            ensure_hash(record)
            self.delay(action_obj)

    def send_from_queue(self, action_objs):
        if not isinstance(action_objs, list):
            action_objs = [action_objs]
        for action_obj in action_objs:
            ensure_hash(action_obj['record'])
        hashes = {action_obj['record']['hash']: action_obj for action_obj in action_objs}
        records = [action_obj['record'] for action_obj in hashes.values()]
        try:
            succeeded, failed = self.action_plugin.send(records, self.content)
        except Exception as err:
            proclog.exception(err)
            succeeded, failed = [], records
        for record in failed:
            hashes[record['hash']]['retry'] -= 1
            if hashes[record['hash']]['delay'] > 0:
                hashes[record['hash']]['delay'] = max(hashes[record['hash']]['freq'], 0) or hashes[record['hash']]['every']
        for record in succeeded:
            hashes[record['hash']]['total'] -= 1
            if hashes[record['hash']]['delay'] > 0:
                hashes[record['hash']]['delay'] = hashes[record['hash']]['every']
        self.update_stats(succeeded, len(succeeded))
        return succeeded, failed

    def send_one(self, loop, action_obj):
        record = action_obj['record']
        success = True
        for _ in range(loop):
            try:
                if self.batch:
                    action_obj.pop('action', '')
                    got_queued = self.plugin.core.mq.send(f"action_{self.uid}", action_obj)
                    if got_queued:
                        action_obj['total'] = 0
                        break
                    else:
                        failed = record
                else:
                    _, failed = self.action_plugin.send(record, self.content)
            except Exception as err:
                proclog.error("Action %s could not be sent: %s", self.name, err, exc_info=err)
                failed = record
            if failed:
                action_obj['retry'] -= 1
                action_obj['delay'] = max(action_obj['freq'], 0) or action_obj['every']
                success = False
                break
            else:
                action_obj['total'] -= 1
                action_obj['delay'] = action_obj['every']
        if success:
            self.update_stats(True, loop)
        return success

    def delay(self, action_obj):
        '''Delay an action'''
        record = action_obj['record']
        if action_obj['total'] == 0 or action_obj['retry'] < 0:
            self.delayed.cleanup(record['hash'], self.uid)
            if action_obj['retry'] < 0:
                self.update_stats(False)
            return
        action_obj['action'] = self
        action_obj['time'] = time.time() + action_obj['delay']
        self.delayed.set_delayed(action_obj)
        proclog.info("Action '%s' will be sent in %d seconds (%d retrie(s) left)",
            self.name, action_obj['delay'], action_obj['retry'])

    def update_stats(self, success, amount = 1):
        '''Update internal action statistics'''
        if amount > 0:
            try:
                if success:
                    self.core.stats.inc('action_success', {'name': self.name}, amount)
                else:
                    self.core.stats.inc('action_error', {'name': self.name}, amount)
            except Exception as err:
                proclog.warning("Error while incrementing stat for action %s",
                    self.name, exc_info=err, extra=dict(plugin='action'))

    def __str__(self):
        return self.name


class DelayedActions(SurvivingThread):
    '''A thread for handling delayed actions'''
    def __init__(self, action, exit_event: Optional[Event] = None):
        exit_event = exit_event or Event()
        super().__init__()
        self.action = action
        self.delayed = {}

        SurvivingThread.__init__(self, exit_event)

    def set_delayed(self, action_obj, update_db=True):
        record_hash = action_obj['record']['hash']
        action_uid = action_obj['action'].uid
        if record_hash not in self.delayed:
            self.delayed[record_hash] = {}
        self.delayed[record_hash][action_uid] = action_obj
        if update_db:
            object_to_write = {
                'action_uid': action_uid,
                'record_hash': record_hash,
                'host': self.action.hostname,
                'every': action_obj['every'],
                'delay': action_obj['delay'],
                'total': action_obj['total'],
                'retry': action_obj['retry'],
                'freq': action_obj['freq'],
            }
            self.action.core.db.write('action.delay', object_to_write, 'action_uid,record_hash')

    def cleanup(self, record_hash, action_uid=None):
        try:
            if record_hash in self.delayed:
                if action_uid:
                    query = ['AND', ['=', 'record_hash', record_hash], ['=', 'action_uid', action_uid], ['=', 'host', self.action.hostname]]
                    self.action.core.db.delete('action.delay', query)
                    if action_uid in self.delayed[record_hash]:
                        del self.delayed[record_hash][action_uid]
                    if not self.delayed[record_hash]:
                        del self.delayed[record_hash]
                else:
                    self.action.core.db.delete('action.delay', ['AND', ['=', 'record_hash', record_hash], ['=', 'host', self.action.hostname]])
                    del self.delayed[record_hash]
        except KeyError as err:
            proclog.exception(err, extra=dict(plugin='action'))

    def send_delayed(self, record_hash, action_uid):
        delayed_record = self.action.core.db.get_one('record', dict(hash=record_hash))
        if delayed_record:
            if delayed_record.get('state') in ['ack', 'close'] or delayed_record.get('snoozed'):
                proclog.info("Record already acked/closed/snoozed. Will not notify")
                self.cleanup(record_hash)
            else:
                try:
                    self.delayed[record_hash][action_uid]['record'] = delayed_record
                    action = self.delayed[record_hash][action_uid]['action']
                    success = action.send_one(1, self.delayed[record_hash][action_uid])
                    if success:
                        self.action.core.db.replace_one('record', {'uid': delayed_record['uid']}, delayed_record)
                    action.delay(self.delayed[record_hash][action_uid])
                except Exception as err:
                    proclog.exception(err)
        else:
            proclog.warning("Record %s does not exist anymore. Will not notify", record_hash)
            self.cleanup(record_hash)

    def start_thread(self):
        while not self.exit.wait(0.1):
            for record_hash in list(self.delayed.keys()):
                for action_uid in list(self.delayed[record_hash].keys()):
                    try:
                        if time.time() >= self.delayed[record_hash][action_uid]['time']:
                            self.send_delayed(record_hash, action_uid)
                    except KeyError:
                        apilog.info("KeyError: Hash %s does not exist anymore", record_hash)
            time.sleep(2)
        apilog.info('Stopped delayed action thread', extra=dict(plugin='action'))

class ActionWorker(Worker):

    def process(self):
        to_ack = []
        for action_obj, msg in self.to_ack:
            record = self.thread.obj.core.db.get_one('record', dict(hash=action_obj['record']['hash']))
            if record:
                action_obj['record'] = record
                if record.get('state') in ['ack', 'close'] or record.get('snoozed'):
                    log.debug("Record %s in batch is already acked, closed or snoozed. Do not notify", record.get('hash'))
                    msg.ack()
                else:
                    to_ack.append((action_obj, msg))
            else:
                to_ack.append((action_obj, msg))
        self.to_ack = to_ack
        if len(to_ack) == 0:
            return
        succeeded, _ = self.thread.obj.send_from_queue([action_obj for action_obj, _ in self.to_ack])
        if succeeded:
            self.thread.obj.core.db.write('record', succeeded, 'hash')
        for action_obj, msg in self.to_ack:
            self.thread.obj.delay(action_obj)
            msg.ack()
        self.to_ack = []
