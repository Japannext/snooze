#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

import inspect

from snooze.utils.config import *

class TestConfig:
    def test_empty(self, tmp_path):
        config = Config(tmp_path)
        assert isinstance(config.core, CoreConfig)
        assert isinstance(config.general, GeneralConfig)
        assert isinstance(config.notifications, NotificationConfig)
        assert config.ldap is None

class TestCoreConfig:
    def test_empty(self, tmp_path):
        config = CoreConfig(tmp_path)
        assert config.listen_addr == '0.0.0.0'
        assert config.port == 5200

    def test_read(self, tmp_path):

        core_path = tmp_path / 'core.yaml'
        data = inspect.cleandoc('''---
        listen_addr: '0.0.0.0'
        port: '5200'
        debug: true
        bootstrap_db: true
        create_root_user: true
        unix_socket: /var/run/snooze/server-test.socket
        no_login: false
        audit_excluded_paths: ['/api/patlite', '/metrics', '/web']
        ssl:
          enabled: true
          certfile: '/etc/pki/tls/certs/snooze.crt'
          keyfile: '/etc/pki/tls/private/snooze.key'
        web:
          enabled: true
          path: /opt/snooze/web
        process_plugins: [rule, aggregaterule, snooze, notification]
        database:
          type: mongo
        ''')
        core_path.write_text(data)

        config = CoreConfig(tmp_path)
        assert config.listen_addr == '0.0.0.0'
        assert config.port == 5200
        assert config.debug == True
        assert config.bootstrap_db == True

class TestHousekeeperConfig:
    def test_empty(self, tmp_path):
        config = HousekeeperConfig(tmp_path)
        assert config

class TestGeneralConfig:
    def test_empty(self, tmp_path):
        config = GeneralConfig(tmp_path)
        assert config

class TestNotificationConfig:
    def test_empty(self, tmp_path):
        config = NotificationConfig(tmp_path)
        assert config

class TestMetadataConfig:
    def test_all_plugins(self):
        metadata_files = SNOOZE_PLUGIN_PATH.glob('*/metadata.yaml')
        plugins = [path.parent.name for path in metadata_files]
        assert plugins
        metadata = {}
        for plugin in plugins:
            metadata[plugin] = MetadataConfig(plugin)

        assert metadata['audit'].name == 'Audit'
        assert metadata['snooze'].name == 'Snooze'
