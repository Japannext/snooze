#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Module for managing loading and writing the configuration files'''

import os
from contextlib import contextmanager
from logging import getLogger
from pathlib import Path
from typing import Optional, List, Any

import yaml
from filelock import FileLock

from snooze import __file__ as SNOOZE_PATH

log = getLogger('snooze.utils.config')

SNOOZE_CONFIG_PATH = Path(os.environ.get('SNOOZE_SERVER_CONFIG', '/etc/snooze/server'))
DEFAULTS = Path(SNOOZE_PATH).parent / 'defaults'

class ReadOnlyConfig:
    '''A class representing a config file at a given path.
    Can only be read.'''
    path: Path
    data: dict
    default_data: dict

    def __init__(self, section: str, default_basedir: Path = DEFAULTS):
        path = SNOOZE_CONFIG_PATH / (section + '.yaml')

        if default_basedir:
            try:
                default_path = default_basedir / path.name
                data = default_path.read_text(encoding='utf-8')
                self.default_data = yaml.safe_load(data)
            except (OSError, yaml.YAMLError):
                self.default_data = {}
        else:
            self.default_data = {}

        self.path = path
        self.data = {}
        # pylint: disable=abstract-class-instantiated
        self.filelock = FileLock(path)
        self.read()

    def read(self):
        '''Read the config file to load the config'''
        data = self.path.read_text(encoding='utf-8')
        self.data = {**self.default_data, **yaml.safe_load(data)}

    def get(self, key: str, default: Optional[Any] = None) -> Any:
        '''Same as data get method'''
        return self.data.get(key, default)

    def __getitem__(self, key: str):
        return self.data[key]

    def dig(self, *keys: List[str], default: Optional[Any] = None) -> Any:
        '''Get a nested key from the config'''
        cursor = self.data
        for key in keys:
            cursor = cursor.get(key)
            if not isinstance(cursor, dict):
                return default
        return cursor

class Config(ReadOnlyConfig):
    '''A class representing a writable config file at a given path.
    Can be explored, and updated with a lock file.'''

    @contextmanager
    def lock(self):
        '''A context manager to filelock the config during an update (lock, read, update, write, release)'''
        self.filelock.acquire()
        self.read()
        try:
            yield # Update the config
        finally:
            self.write()
            self.filelock.release()

    def write(self):
        '''Write a new config to the config file'''
        data = yaml.dump(self.data)
        self.path.write_text(data, encoding='utf-8')

    def __setitem__(self, key: str, value: Any):
        self.set(key, value)

    def set(self, key: str, value: Any):
        '''Rewrite a config key with a given value'''
        with self.lock():
            self.data[key] = value

class Metadata(ReadOnlyConfig):
    '''A class to fetch metadata configuration'''
    def __init__(self, plugin_name: str):
        path = Path(SNOOZE_PATH).parent / 'plugins/core' / plugin_name / 'metadata.yaml'
        if path.is_file():
            ReadOnlyConfig.__init__(self, path, False)
        else:
            self.default_data = {}
            self.data = {}
