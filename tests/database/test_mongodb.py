'''Test for the Mongodb database endpoint'''

from unittest.mock import patch
from contextlib import contextmanager

import pytest
import mongomock
from pymongo import MongoClient
from bson.binary import UuidRepresentation
from bson.codec_options import CodecOptions

from snooze.utils.condition import *
from snooze.database.mongodb import *
from snooze.utils.model import *

@pytest.fixture(scope='function')
def db():
    client = mongomock.MongoClient()
    client.drop_database('snooze')
    return client['snooze']

@contextmanager
def transaction_mock(_client, **_kwargs):
    '''Mongomock doesn't support sessions yet, so we're mocking our
    transaction to return session=None. This allow mongodb calls
    to proceed as if there was no session at all. Thus, we cannot
    test race conditions'''
    yield None

def get_collection_mock(database, collection):
    return database.get_collection(collection)

@patch('snooze.database.mongodb.transaction', transaction_mock)
@patch('snooze.database.mongodb.get_collection', get_collection_mock)
class TestEndpoint:
    def test_init(self, db):
        print(db)
        endpoint = Endpoint(db, Rule)
        assert endpoint

    def test_search(self, db):
        endpoint = Endpoint(db, Rule)
        results = endpoint.search()
        assert results == []

    def test_create(self, db):
        endpoint = Endpoint(db, Rule)
        user = SnoozeUser(name='test', method='local')
        rules = [
            Rule(name='Rule 01', condition=AlwaysTrue(), snooze_user=user),
            Rule(name='Rule 02', condition=Equals(field='process', value='systemd'), snooze_user=user)
        ]
        created = endpoint.create(rules)
        assert created
        results = endpoint.search()
        assert len(results) == 2
        assert results[0].name == 'Rule 01'
        assert results[0].uid == rules[0].uid

        assert results[1].name == 'Rule 02'
        assert results[1].uid == rules[1].uid

    def test_get_by_uid(self, db):
        endpoint = Endpoint(db, Rule)
        user = SnoozeUser(name='test', method='local')
        rules = [
            Rule(name='Rule 01', condition=AlwaysTrue(), snooze_user=user),
            Rule(name='Rule 02', condition=Equals(field='process', value='systemd'), snooze_user=user)
        ]
        created = endpoint.create(rules)
        rule1 = endpoint.get_by_uid(rules[0].uid)
        assert rule1 is not None
        assert rule1.uid == rules[0].uid
        assert rule1.name == rules[0].name
        assert rule1.condition == rules[0].condition
        assert rule1.snooze_user == rules[0].snooze_user

    def test_replace(self, db):
        endpoint = Endpoint(db, Rule)
        user = SnoozeUser(name='test', method='local')
        rules = [
            Rule(name='Rule 01', condition=AlwaysTrue(), snooze_user=user),
            Rule(name='Rule 02', condition=Equals(field='process', value='systemd'), snooze_user=user, comment='my comment')
        ]
        print('Rules: ', rules)
        created = endpoint.create(rules)
        print('Created: ', created)
        assert created
        new_rule = Rule(
            name='Rule 03',
            condition=Equals(field='process', value='rsyslog'),
            snooze_user=SnoozeUser(name='test2', method='local'),
        )
        endpoint.replace(rules[1].uid, new_rule)
        print('Search: ', endpoint.search())
        result = endpoint.get_by_uid(rules[1].uid)
        assert result is not None
        assert result.name == new_rule.name
        assert result.condition == new_rule.condition
        assert result.snooze_user == new_rule.snooze_user
        assert result.comment == ''

    def test_patch_one(self, db):
        endpoint = Endpoint(db, Rule)
        user = SnoozeUser(name='test', method='local')
        rules = [
            Rule(name='Rule 01', condition=AlwaysTrue(), snooze_user=user),
            Rule(name='Rule 02', condition=Equals(field='process', value='systemd'), snooze_user=user)
        ]
        created = endpoint.create(rules)
        new_rule = Partial[Rule](
            condition=Equals(field='process', value='rsyslog'),
        )
        print('db.rule.find(<uid>): ', db.rule.find({'uid': rules[0].uid})[0])
        updated = endpoint.patch_one(rules[0].uid, new_rule)
        print(updated)
        result = endpoint.get_by_uid(rules[0].uid)
        assert result.name == 'Rule 01'
        assert result.condition == Equals(field='process', value='rsyslog')
        assert result.snooze_user == user

    def test_patch_many(self, db):
        endpoint = Endpoint(db, Rule)
        user = SnoozeUser(name='test', method='local')
        rules = [
            Rule(name='Rule 01', condition=AlwaysTrue(), snooze_user=user, comment='My comment 1'),
            Rule(name='Rule 02', condition=Equals(field='process', value='systemd'), snooze_user=user)
        ]
        created = endpoint.create(rules)
        new_rule = Partial[Rule](
            comment='New comment'
        )
        uids = [rule.uid for rule in rules]
        updated = endpoint.patch_many(uids, new_rule)
        results = endpoint.search()
        for result in results:
            assert result.comment == 'New comment'

    def test_delete_one(self, db):
        ...

    def test_delete_many(self, db):
        ...
