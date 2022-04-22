#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#
'''Testing comment routes'''

class TestCommentRoute:

    data = {
        'record': [
            {'name': 'r1'},
            {'name': 'r2'},
            {'name': 'r3'},
        ],
    }

    def test_post(self, client):
        record1 = client.simulate_get('/api/record', params=dict(s=['=', 'name', 'r1'])).json['data'][0]
        comment1 = {'record_uid': record1['uid'], 'message': 'comment 1'}
        resp = client.simulate_post('/api/comment', json=[comment1])
        assert resp.status == '201 Created'
        added = resp.json['data']['added']
        assert len(added) == 1

        record1_after = client.simulate_get('/api/record', params=dict(s=['=', 'name', 'r1'])).json['data'][0]
        assert record1_after['comment_count'] == 1

        comment2 = {'record_uid': record1['uid'], 'message': 'comment 2'}
        resp = client.simulate_post('/api/comment', json=[comment2])
        assert resp.status == '201 Created'
        added = resp.json['data']['added']
        assert len(added) == 1

        record1_after_bis = client.simulate_get('/api/record', params=dict(s=['=', 'name', 'r1'])).json['data'][0]
        assert record1_after_bis['comment_count'] == 2

    def test_put(self, client):
        record1 = client.simulate_get('/api/record', params=dict(s=['=', 'name', 'r1'])).json['data'][0]
        comment1 = {'record_uid': record1['uid'], 'message': 'comment 1'}
        resp = client.simulate_post('/api/comment', json=[comment1])
        assert resp.status == '201 Created'
        added = resp.json['data']['added']
        assert len(added) == 1
        comment_uid = added[0]['uid']

        added[0]['message'] = 'message2'
        resp = client.simulate_put('/api/comment', json=[added[0]])
        assert resp.status == '201 Created'
        updated = resp.json['data']['updated']
        assert len(updated) == 1

        record1_after = client.simulate_get('/api/record', params=dict(s=['=', 'name', 'r1'])).json['data'][0]
        assert record1_after['comment_count'] == 1
