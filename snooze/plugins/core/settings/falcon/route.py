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
from snooze.utils.config import GeneralConfig, NotificationConfig, LdapConfig, HousekeeperConfig
from snooze.utils.typing import SettingUpdatePayload

log = getLogger('snooze.api')

CONFIGS = {
    'general': GeneralConfig,
    'notifications': NotificationConfig,
    'ldap': LdapConfig,
    'housekeeper': HousekeeperConfig,
}

class SettingsRoute(BasicRoute):
    @authorize
    def on_get(self, req, resp, conf=''):
        '''Fetch a config file data.
        Secrets are protected thanks to the `Field(exclude=True)` of pydantic.
        ValidationError are server side errors (the local config file is broken)'''
        section = req.params.get('c') or conf
        checksum = req.params.get('checksum')

        log.debug("Loading config file %s", section)
        resp.content_type = falcon.MEDIA_JSON
        try:
            result_dict = CONFIGS[section](basedir=self.api.core.basedir).dict()
        except KeyError:
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
            config_class = CONFIGS[section]
            config_class().update(payload.conf)
            self.api.core.reload(section)
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
