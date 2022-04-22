#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#
'''Testing reload route'''

class TestReloadPluginRoute:
    def test_reload(self, client):
        resp = client.simulate_post('/api/reload/snooze')
        assert resp.status == '200 OK'

    def test_reload_propagate(self, client):
        resp = client.simulate_post('/api/reload/snooze?propagate')
        assert resp.status == '202 Accepted'
