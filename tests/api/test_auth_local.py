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
from falcon.testing import TestClient

from snooze.api import Api
from snooze.core import Core

log = getLogger('snooze.tests.api')

"""
def test_basic_auth(core):
    users = [{"name": "root", "method": "local", "enabled": True}]
    core.db.write('user', users)
    user_passwords = [{"name": "root", "method": "local", "password": sha256("root".encode('utf-8')).hexdigest()}]
    core.db.write('user.password', user_passwords)
    token = str(b64encode("{}:{}".format('root', 'root').encode('utf-8')), 'utf-8')
    headers = {'Authorization': 'Basic {}'.format(token)}
    log.debug(headers)
    client = testing.TestClient(core.api.handler, headers=headers)
    log.debug('Attempting Basic auth')
    result = client.simulate_post('/api/login/local').json
    log.debug("Received {}".format(result))
    assert result
    assert result['token']
"""

class TestLocalAuthRoute:
    data = {
        'user': [
            {'name': 'test', 'method': 'local', 'enabled': True},
        ],
    }

    def test_basic_auth(self, core):
        credentials = b64encode('test:secret123'.encode('utf-8')).decode('utf-8')
        headers = {'Authorization': f"Basic {credentials}"}
        password_hash = sha256('secret123'.encode('utf-8')).hexdigest()
        user_password = {'name': 'test', 'method': 'local', 'password': password_hash}
        core.db.write('user.password', [user_password])
        client = TestClient(core.api.handler, headers=headers)
        resp = client.simulate_post('/api/login/local')
        assert resp.status == '200 OK'
        token = resp.json['token']
