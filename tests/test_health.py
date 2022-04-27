'''Test the health route'''

from snooze.health import *

def test_thread_status(started_core):
    status = {'issues': []}
    healths = thread_status(started_core, status)
    print(healths)
    assert all(health == 'ok' for health in healths)
    assert status['threads']['tcp'] == {'alive': True}
    assert status['threads']['housekeeper'] == {'alive': True}
    assert status['threads']['cluster'] == {'alive': True}
    assert status['threads']['tcp'] == {'alive': True}
    assert status['issues'] == []

def test_mq_status(core):
    status = {'issues': []}
    healths = mq_status(core.mq, status)
    assert all(health == 'ok' for health in healths)
    for thread in status['mq']['threads'].values():
        assert thread == {'alive': True}

class TestHealthRoute:
    def test_get(self, started_client):
        started_client.simulate_get('/api/health')
