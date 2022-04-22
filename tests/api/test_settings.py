#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#
'''Testing settings route'''

class TestSettingsRoute:
    def test_get(self, client):
        resp = client.simulate_get('/api/settings/general')
        assert resp.status == '200 OK'
        data = resp.json['data']
        assert data['default_auth_backend'] == 'local'

    def test_put(self, client):
        data = {'metric_enabled': False}
        resp = client.simulate_put('/api/settings/general', json=data)
        assert resp.status == '200 OK'

    def test_put_propagate(self, client):
        data = {'metric_enabled': False}
        resp = client.simulate_put('/api/settings/general?propagate', json=data)
        assert resp.status == '202 Accepted'
