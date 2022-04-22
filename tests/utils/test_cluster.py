#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

import socket
from unittest.mock import patch

import pytest
import responses

from snooze.utils.cluster import *
from snooze.utils.exceptions import *
from snooze.utils.typing import *

def resolver(hostname):
    if hostname == 'host01':
        return '10.0.0.11'
    if hostname == 'host02':
        return '10.0.0.12'
    if hostname == 'host03':
        return '10.0.0.13'
    raise socket.gaierror('[Errno -2] Name or service not known')

def mock_ifaddresses(interface):
    if interface == 'eth0':
        return {
            17: [{'addr': 'fa:gg:gg:gg:gg:gg', 'broadcast': 'ff:ff:ff:ff:ff:ff'}],
            2: [{'addr': '10.0.0.11', 'netmask': '255.255.255.0', 'broadcast': '10.0.0.255'}],
            10: [{'addr': 'fe80:ffff:ffff:ffff:ffff%eth0', 'netmask': 'ffff:ffff:ffff:ffff::/64'}],
        }
    if interface == 'lo':
        return {
            17: [{'addr': '00:00:00:00:00:00', 'peer': '00:00:00:00:00:00'}],
            2: [{'addr': '127.0.0.1', 'netmask': '255.0.0.0', 'broadcast': '127.0.0.1'}],
            10: [{'addr': '::1', 'netmask': 'ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128'}],
        }
    else:
        raise ValueError('You must specify a valid interface name.')

def mock_who_am_i(members):
    return members[0], members[1:]

class TestStandaloneCluster:
    def test_status(self, config):
        hostname = socket.gethostname()
        cluster = Cluster(config.core)
        status = cluster.status()
        assert isinstance(status, PeerStatus)
        assert status.host == 'localhost'
        assert status.port == 5200
        assert status.healthy == True

    def test_member_status(self, config):
        hostname = socket.gethostname()
        cluster = Cluster(config.core)
        statuses = cluster.members_status()
        assert len(statuses) == 1
        assert statuses[0].host == 'localhost'
        assert statuses[0].port == 5200
        assert statuses[0].healthy == True

class TestRealCluster:

    configs = {
        'core': {
            'cluster': {
                'enabled': True,
                'members': [
                    {'host': 'host01'},
                    {'host': 'host02'},
                    {'host': 'host03'},
                ],
            },
            'backup': {'enabled': False},
        },
    }

    @patch('snooze.utils.cluster.who_am_i', mock_who_am_i)
    def test_status(self, config):
        hostname = socket.gethostname()
        cluster = Cluster(config.core)
        status = cluster.status()
        assert isinstance(status, PeerStatus)
        assert status.host == 'host01'
        assert status.port == 5200
        assert status.healthy == True

    @responses.activate
    @patch('snooze.utils.cluster.who_am_i', mock_who_am_i)
    def test_member_status(self, config):
        host02_status = {'host': 'host02', 'port': 5200, 'version': '1.x.x', 'healthy': True}
        host03_status = {'host': 'host03', 'port': 5200, 'version': '1.x.x', 'healthy': True}
        responses.add(responses.GET, 'https://host02:5200/api/cluster', status=200, json=host02_status)
        responses.add(responses.GET, 'https://host03:5200/api/cluster', status=200, json=host03_status)
        cluster = Cluster(config.core)
        statuses = cluster.members_status()
        assert len(statuses) == 3
        assert statuses[0].host == 'host01'
        assert statuses[0].port == 5200
        assert statuses[0].healthy == True
        assert statuses[1].host == 'host02'
        assert statuses[1].port == 5200
        assert statuses[1].version == '1.x.x'
        assert statuses[1].healthy == True
        assert statuses[2].host == 'host03'
        assert statuses[2].port == 5200
        assert statuses[2].version == '1.x.x'
        assert statuses[2].healthy == True

@patch('snooze.utils.cluster.interfaces', lambda: ['lo', 'eth0'])
@patch('snooze.utils.cluster.ifaddresses', mock_ifaddresses)
def test_who_am_i():
    members = [
        HostPort(host='10.0.0.11'),
        HostPort(host='10.0.0.12'),
        HostPort(host='10.0.0.13'),
    ]
    myself, others = who_am_i(members)
    assert myself.host == '10.0.0.11'
    assert len(others) == 2
    assert others[0].host == '10.0.0.12'
    assert others[1].host == '10.0.0.13'

@patch('snooze.utils.cluster.interfaces', lambda: ['lo', 'eth0'])
@patch('snooze.utils.cluster.ifaddresses', mock_ifaddresses)
def test_who_am_i_not_in_cluster():
    members = [
        HostPort(host='10.0.0.12'),
        HostPort(host='10.0.0.13'),
        HostPort(host='10.0.0.14'),
    ]
    with pytest.raises(SelfNotInCluster):
        myself, others = who_am_i(members)

@patch('snooze.utils.cluster.interfaces', lambda: ['lo', 'eth0'])
@patch('snooze.utils.cluster.ifaddresses', mock_ifaddresses)
def test_who_am_i_too_much_in_cluster():
    members = [
        HostPort(host='10.0.0.11'),
        HostPort(host='10.0.0.11'),
        HostPort(host='10.0.0.12'),
        HostPort(host='10.0.0.13'),
    ]
    with pytest.raises(SelfTooMuchInCluster):
        myself, others = who_am_i(members)

@patch('snooze.utils.cluster.interfaces', lambda: ['lo', 'eth0'])
@patch('snooze.utils.cluster.ifaddresses', mock_ifaddresses)
@patch('socket.gethostbyname', resolver)
def test_who_am_i_dns():
    members = [
        HostPort(host='host01'),
        HostPort(host='host02'),
        HostPort(host='host03'),
    ]
    myself, others = who_am_i(members)
    assert myself.host == 'host01'
    assert len(others) == 2
    assert others[0].host == 'host02'
    assert others[1].host == 'host03'

@patch('snooze.utils.cluster.interfaces', lambda: ['lo', 'eth0'])
@patch('snooze.utils.cluster.ifaddresses', mock_ifaddresses)
@patch('socket.gethostbyname', resolver)
def test_who_am_i_dns_not_resolvable():
    members = [
        HostPort(host='host01'),
        HostPort(host='host02'),
        HostPort(host='host04'),
    ]
    with pytest.raises(NonResolvableHost) as exc_info:
        myself, others = who_am_i(members)
    assert exc_info.value.host == 'host04'
