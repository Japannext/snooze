#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

class TestCore:
    data = {'record': []}

    def test_load_plugins(self, core):
        plugins = [plugin.name for plugin in core.plugins]
        default_plugins = ['record', 'rule', 'aggregaterule', 'snooze', 'notification']
        for plugin in default_plugins:
            assert plugin in plugins

    def test_process_record(self, core):
        record = {'a': '1', 'b': '2'}
        core.process_record(record)
        search = core.db.search('record', ['AND', ['=', 'a', '1'], ['=', 'b', '2']])
        assert all(plugin in search['data'][0]['plugins'] for plugin in ['rule', 'aggregaterule', 'snooze', 'notification'])

    def test_process_ok(self, core):
        records = core.db.search('record')
        print(records)
        core.config.general.ok_severities = ['ok']
        record = {'severity': 'ok'}
        core.process_record(record)
        data = core.db.search('record')['data'][0]
        assert data['state'] == 'close'
