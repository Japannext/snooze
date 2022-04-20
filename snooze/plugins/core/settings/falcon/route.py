#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

import hashlib
from logging import getLogger

import falcon
from pydantic import ValidationError

from snooze.api.routes import BasicRoute
from snooze.utils.functions import authorize
from snooze.utils.typing import SettingUpdatePayload

log = getLogger('snooze.api')

class SettingsRoute(BasicRoute):
    @authorize
    def on_get(self, req, resp, section):
        '''Fetch a config file data.
        Secrets are protected thanks to the `Field(exclude=True)` of pydantic.
        ValidationError are server side errors (the local config file is broken)'''

        log.debug("Loading config file %s", section)
        resp.content_type = falcon.MEDIA_JSON
        try:
            config = getattr(self.core.config, section)
            checksum = hashlib.md5(config.json()).encode('utf-8').hexdigest()
        except AttributeError:
            resp.status = falcon.HTTP_400
            resp.media = {'message': f"Unknown config '{section}'"}
        except ValidationError as err:
            resp.status = falcon.HTTP_503
            resp.media = {'message': str(err)}
        if result_dict:
            result_dict = {k:v for k,v in result_dict.items() if 'password' not in k}
            dict_checksum = hashlib.md5(repr([result_dict]).encode('utf-8')).hexdigest()
            if checksum != dict_checksum:
                result = {'data': [result_dict], 'count': 1, 'checksum': dict_checksum}
            else:
                result = {'count': 0}
            resp.media = result
            if 'error' in result_dict.keys():
                resp.status = falcon.HTTP_503
            else:
                resp.status = falcon.HTTP_200
        else:
            resp.media = {}
            resp.status = falcon.HTTP_404

    @authorize
    def on_put(self, req, resp, conf=''):
        '''Action trigerred when a user change a setting section on the web interface'''
        '''Rewrite a config file on the server.
        ValidationError are client side.
        A refresh of the class is needed (`self.api.core.reload(section)`).
        Sharing within the cluster is needed.
        '''
        section = req.params.get('c') or conf
        resp.content_type = falcon.MEDIA_JSON

        log.debug("Trying write to configfile %s: %s", section, req.media)
        try:
            payload = SettingUpdatePayload(**req.media[0])
            self.core.setting_update(payload.section, payload.data)
            self.core.sync_setting_update()
            resp.status = falcon.HTTP_CREATED
            resp.media = {'data': f"Config {section} reloaded"}
        except KeyError:
            resp.status = falcon.HTTP_400
            resp.media = {'data': f"Unknown config '{section}'"}
        except ValidationError as err:
            resp.status = falcon.HTTP_400
            resp.media = {'data': str(err)}
        except Exception as err:
            log.exception(err)
            resp.status = falcon.HTTP_503
            resp.media = {'data': str(err)}

class SettingsUpdateRoute(BasicRoute):
    '''A falcon route to reload one's token'''
    auth = {'auth_disabled': True}

    def on_put(self, req, resp):
        '''Action triggerred by a cluster member to propagate a setting change'''

        try:
            payload = SettingUpdatePayload(**req.media)
        except ValidationError as err: # Bad payload
            raise err
        try:
            self.core.setting_update(payload.section, payload.data)
        except ValidationError as err: # Bad data for config
            raise err
