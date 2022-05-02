#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

from dateutil import parser
from freezegun import freeze_time

from snooze.utils.time_constraints import *

class TestLogic:

    def test_or_match(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00')
        c1 = DatetimeConstraint(until='2021-07-01T14:30:00+09:00')
        c2 = WeekdaysConstraint(week={1: True, 2: True, 3: True, 4: True})
        c3 = TimeConstraint(time1='11:00+09:00', time2='15:00+09:00')
        c = And(constraints=[c1, c2, c3])
        assert c.match(record_date) == True

    def test_or_miss(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00')
        c1 = DatetimeConstraint(until='2021-07-01T14:30:00+09:00')
        c2 = WeekdaysConstraint(week={6: True, 0: True})
        c = And(constraints=[c1, c2])
        assert c.match(record_date) == False

    def test_match_any_same_type(self):
        record_date = parser.parse('2021-07-01T23:00:00+09:00')
        c1 = TimeConstraint(time1='00:00+09:00', time2='02:00+09:00')
        c2 = TimeConstraint(time1='22:00+09:00', time2='23:59+09:00')
        c = Or(constraints=[c1, c2])
        assert c.match(record_date) == True

class TestDatetimeConstraint:

    def test_until_true(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00')
        tc = DatetimeConstraint(until='2021-07-01T14:30:00+09:00')
        assert tc.match(record_date) == True

    def test_until_false(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00')
        tc = DatetimeConstraint(until='2021-07-01T11:30:00+09:00')
        assert tc.match(record_date) == False

    def test_from_true(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00')
        tc = DatetimeConstraint(date_from='2021-07-01T10:30:00+09:00')
        assert tc.match(record_date) == True

    def test_from_false(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00')
        tc = DatetimeConstraint(date_from='2021-07-01T12:30:00+09:00')
        assert tc.match(record_date) == False

class TestWeekdaysConstraint:

    def test_weekday_true(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00') # Thursday (4's day of the week)
        tc = WeekdaysConstraint(week={4: True})
        assert tc.match(record_date) == True

    def test_weekday_false(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00') # Thursday (4's day of the week)
        tc = WeekdaysConstraint(week={6: True, 0: True})
        assert tc.match(record_date) == False

class TestTimeConstraint:

    def test_from_true(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00')
        tc = TimeConstraint(time1='10:00+09:00')
        assert tc.match(record_date) == True
        tc = TimeConstraint(time1='12:00+09:00')
        assert tc.match(record_date) == True

    def test_from_false(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00')
        tc = TimeConstraint(time1='14:00+09:00')
        assert tc.match(record_date) == False

    def test_until_true(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00')
        tc = TimeConstraint(until='14:00+09:00')
        assert tc.match(record_date) == True
        tc = TimeConstraint(until='12:00+09:00')
        assert tc.match(record_date) == True

    def test_until_false(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00')
        tc = TimeConstraint(until='10:00+09:00')
        assert tc.match(record_date) == False

    def test_range_true(self):
        record_date = parser.parse('2021-07-01T12:00:00+09:00')
        tc = TimeConstraint(time1='10:00+09:00', time2='14:00+09:00')
        assert tc.match(record_date) == True

    def test_range_false(self):
        record_date = parser.parse('2021-07-01T08:00:00+09:00')
        tc = TimeConstraint(time1='10:00+09:00', time2='14:00+09:00')
        assert tc.match(record_date) == False

    def test_over_midnight(self):
        record_date = parser.parse('2021-07-01T01:00:00+09:00')
        tc = TimeConstraint(time1='23:00+09:00', time2='02:00+09:00')
        assert tc.match(record_date) == True

    def test_over_midnight_miss(self):
        record_date = parser.parse('2021-07-01T03:00:00+09:00')
        tc = TimeConstraint(time1='23:00+09:00', time2='02:00+09:00')
        assert tc.match(record_date) == False
