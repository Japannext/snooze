#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Module for Local user authentication'''

from base64 import b64decode
from hashlib import sha256
from logging import getLogger

import falcon

from snooze.api.auth import AuthRoute

log = getLogger('snooze.auth.local')

