#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#
'''Test ClusterRoute'''

class TestClusterRoute:

    def test_status(self, client):
        resp = client.simulate_get('/api/cluster')
        assert resp.status == '200 OK'
        data = resp.json['data']
        assert len(data) == 1
        assert data[0]['host'] == 'localhost'
        assert data[0]['port'] == 5200
