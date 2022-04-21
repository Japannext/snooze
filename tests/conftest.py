#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Configuration and fixtures for testing'''

from logging import getLogger

import mongomock
import pytest
import yaml
from falcon.testing import TestClient
from pytest_data.functions import get_data

from snooze.core import Core, MAIN_THREADS
from snooze.db.database import Database
from snooze.utils.config import Config

log = getLogger('tests')

DEFAULT_CONFIG = {
    'core': {
        'database': {'type': 'mongo', 'host': 'localhost', 'port': 27017},
        'socket_path': './test.socket',
        'init_sleep': 0,
        'stats': False,
        'backup': {'enabled': False},
    },
}

def write_data(database, request):
    '''Write data fetch with get_data(request, 'data') to a database'''
    data = get_data(request, 'data', {})
    for key, value in data.items():
        database.delete(key, [], True)
    for key, value in data.items():
        database.write(key, value)

@pytest.fixture(name='config', scope='function')
def fixture_config(tmp_path, request):
    '''Fixture for writable configuration files returning a Config'''
    configs = {**DEFAULT_CONFIG, **get_data(request, 'configs', {})}

    for section, data in configs.items():
        path = tmp_path / f"{section}.yaml"
        path.write_text(yaml.dump(data), encoding='utf-8')

    return Config(tmp_path)

@pytest.fixture(name='db', scope='function')
@mongomock.patch('mongodb://localhost:27017')
def fixture_db(config, request):
    '''Fixture returning a mocked mongodb Database'''
    database = Database(config.core.database)
    write_data(database, request)
    return database

@pytest.fixture(name='core', scope='function')
@mongomock.patch('mongodb://localhost:27017')
def fixture_core(config, request):
    '''Fixture returning a Core'''
    allowed_threads = get_data(request, 'allowed_threads') or MAIN_THREADS
    core = Core(config.basedir, allowed_threads)
    write_data(core.db, request)
    return core

@pytest.fixture(name='api', scope='function')
def fixture_api(core):
    '''Fixture returning an Api'''
    return core.api

@pytest.fixture(name='client', scope='function')
def fixture_client(api):
    '''Fixture returning a falcon TestClient'''
    token = api.get_root_token()
    log.info("Token obtained from get_root_token: %s", token)
    headers = {'Authorization': f"JWT {token}"}
    return TestClient(api.handler, headers=headers)
