'''Test the base API'''

import pytest
from falcon.testing import TestClient

from snooze.utils.model import *

class TestBaseApiRoute:

    data = {
        'rule': [
            {'name': 'rule 1'},
            {'name': 'rule 2'},
        ],
    }

    def test_on_get_search_everything(self, client):
        resp = client.simulate_get('/api/rule')
        assert resp.status == '200 OK'
        assert len(resp.json) == 2

    def test_on_get_uid(self, client):
        uid = ''
        resp = client.simulate_get(f"/api/rule/{uid}")
        assert resp.status == '200 OK'
        assert resp.json['name'] == ''

    def test_on_post(self, client):
        data = {'name': 'rule 03', 'condition': {'type': '=', 'field': 'a', 'value': '1'}}
        resp = client.simulate_post('/api/rule', json=data)
        assert resp.status == '202 CREATED'
        assert isinstance(resp.json['uid'], str)

    def test_on_put(self, client):
        ...

    def test_on_patch(self, client):
        ...

    def test_on_delete(self, client):
        ...

