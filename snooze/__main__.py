#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Initialize the different threads'''

import logging.config
import os
import sys

from logging import getLogger

from snooze.core import Core
from snooze.api import Api
from snooze.utils import config

def setup_logging(conf):
    '''Initialize the python logger'''
    logging_config = config('logging')
    if os.environ.get('SNOOZE_DEBUG', conf.get('debug', False)):
        try:
            logging_config['handlers']['console']['level'] = 'DEBUG'
        except KeyError:
            pass
        try:
            logging_config['handlers']['file']['level'] = 'DEBUG'
        except KeyError:
            pass
        try:
            logging_config['loggers']['snooze']['level'] = 'DEBUG'
        except KeyError:
            pass
    logging.config.dictConfig(logging_config)
    log = getLogger('snooze')
    log.debug("Log system on")
    return log

def app(conf=None):
    '''Used to initialize the application in Docker Heroku'''
    if conf is None:
        conf = {}
    conf.update(config())
    setup_logging(conf)
    core = Core(conf)

    api = Api(core)
    return api.handler

def main(conf=None):
    '''Main thread when running snooze-server executable'''
    log = setup_logging(conf)
    core = Core()

    try:
        core.start()
        if core.exit_event.wait():
            core.stop()
            sys.exit(0)
    except (KeyboardInterrupt, SystemExit):
        core.stop()
        sys.exit(1)

if __name__ == '__main__':
    main()
