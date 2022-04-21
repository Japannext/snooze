#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#


import json
import mongomock
from base64 import b64encode
from hashlib import sha256
from logging import getLogger


import pytest
import yaml
from falcon import testing

from snooze.api import Api
from snooze.core import Core

log = getLogger('snooze.tests.api')

@mongomock.patch('mongodb://localhost:27017')
def test_basic_auth(core):
    api = Api(core)
    users = [{"name": "root", "method": "local", "enabled": True}]
    core.db.write('user', users)
    user_passwords = [{"name": "root", "method": "local", "password": sha256("root".encode('utf-8')).hexdigest()}]
    core.db.write('user.password', user_passwords)
    token = str(b64encode("{}:{}".format('root', 'root').encode('utf-8')), 'utf-8')
    headers = {'Authorization': 'Basic {}'.format(token)}
    log.debug(headers)
    client = testing.TestClient(api.handler, headers=headers)
    log.debug('Attempting Basic auth')
    result = client.simulate_post('/api/login/local').json
    log.debug("Received {}".format(result))
    assert result
    assert result['token']
