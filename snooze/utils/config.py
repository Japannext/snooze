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
from datetime import datetime, timedelta
from typing import Optional, List, Any, Dict, Literal

import yaml
from filelock import FileLock
from pydantic import BaseModel, Field, PrivateAttr, validator

from snooze import __file__ as SNOOZE_PATH
from snooze.utils.typing import RouteArgs

log = getLogger('snooze.utils.config')

SNOOZE_CONFIG = Path(os.environ.get('SNOOZE_SERVER_CONFIG', '/etc/snooze/server'))
SNOOZE_PLUGIN_PATH = Path(SNOOZE_PATH).parent / 'plugins/core'
DEFAULTS = Path(SNOOZE_PATH).parent / 'defaults'

class ReadOnlyConfig(BaseModel):
    '''A class representing a config file at a given path.
    Can only be read.'''
    _path: Path = PrivateAttr()

    def __init__(self):
        data = self._read()
        super().__init__(**data)

    def _read(self) -> dict:
        '''Read the config file and return the raw dict'''
        return yaml.safe_load(self._path.read_text(encoding='utf-8'))

    def refresh(self):
        '''Read the config file to load the config'''
        data = self._read()
        for key, value in data.items():
            setattr(self, key, value)

    def __getitem__(self, key: str):
        return getattr(self, key)

    def dig(self, *keys: List[str], default: Optional[Any] = None) -> Any:
        '''Get a nested key from the config'''
        try:
            cursor = getattr(self, keys[0])
            for key in keys[1:]:
                cursor = getattr(cursor, key)
            return cursor
        except AttributeError:
            return default

class WritableConfig(ReadOnlyConfig):
    '''A class representing a writable config file at a given path.
    Can be explored, and updated with a lock file.'''
    _filelock: FileLock = PrivateAttr()

    def __init__(self):
        self._filelock = FileLock(self._path, timeout=1)
        super().__init__()

    @contextmanager
    def _lock(self):
        '''A context manager to filelock the config during an update (lock, read, update, write, release)'''
        self._filelock.acquire()
        self.refresh()
        try:
            yield # Update the config
        finally:
            self._update()
            self._filelock.release()

    def _update(self):
        '''Write a new config to the config file'''
        data = yaml.dump(self.dict())
        self._path.write_text(data, encoding='utf-8')

    def __setitem__(self, key: str, value: Any):
        self._set(key, value)

    def __setattr__(self, key: str, value: Any):
        self._set(key, value)

    def _set(self, key: str, value: Any):
        '''Rewrite a config key with a given value'''
        with self._lock():
            setattr(self, key, value)

    def update(self, values: dict):
        '''Update the config with a dictionary'''
        with self._lock():
            for key, value in values.items():
                setattr(self, key, value)

class MetadataConfig(ReadOnlyConfig):
    '''A class to fetch metadata configuration'''
    class_name: str = Field(alias='class')
    auto_reload: bool = False
    default_sorting: Optional[str] = None
    default_ordering: bool = True
    widgets: dict = Field(default_factory=dict)
    action_form: dict  = Field(default_factory=dict)
    audit: bool = True
    provides: List[str] = Field(default_factory=list)
    icon: str = 'question-circle'
    routes: Dict[str, RouteArgs]
    route_defaults: RouteArgs
    options: dict = Field(default_factory=dict)
    search_fields: List[str] = Field(default_factory=list)

    def __init__(self, plugin_name: str):
        self._path = SNOOZE_PLUGIN_PATH / plugin_name / 'metadata.yaml'
        data = self._read()
        super().__init__(**data)

class LdapConfig(WritableConfig):
    '''A dataclass representing the LDAP configuration'''
    _path = SNOOZE_CONFIG / 'ldap_auth.yaml'

    enabled: bool
    base_dn: str
    user_filter: str
    bind_dn: str
    bind_password: str = Field(exclude=True)
    host: str
    port: int = 636
    group_dn: Optional[str] = None
    email_attribute: str = 'mail'
    display_name_attribute: str = 'cn'
    member_attribute: str = 'memberof'

class SslConfig(BaseModel):
    '''SSL configuration'''
    enabled: bool = True
    certfile: Path = Field(env='SNOOZE_CERT_FILE')
    keyfile: Path = Field(env='SNOOZE_KEY_FILE')

class CoreConfig(ReadOnlyConfig):
    '''pydantic model representing the core config'''
    _path = SNOOZE_CONFIG / 'core.yaml'

    listen_addr: str = '0.0.0.0'
    port: int = 5200
    ssl: SslConfig
    debug: bool = False
    bootstrap_db: bool = True
    unix_socket: Optional[Path] = Path('/var/run/snooze/server.socket')
    no_login: bool = False
    audit_excluded_paths: List[str] = ('/api/patlite', '/metrics', '/web')
    ssl: dict
    process_plugins: List[str] = ('rule', 'aggregaterule', 'snooze', 'notification')
    database: dict
    init_sleep: int = 5
    create_root_user: bool = False

class GeneralConfig(WritableConfig):
    '''pydantic model representing the general config'''
    _path = SNOOZE_CONFIG / 'general.yaml'

    default_auth_backend: Literal['local', 'ldap'] = 'local'
    local_users_enabled: bool = True
    metrics_enabled: bool = True
    anonymous_enabled: bool = False
    ok_severities: List[str] = ('ok', 'success')

    @validator('ok_severities', each_item=True)
    def normalize_severities(cls, value): # pylint: disable=no-self-argument,no-self-use
        '''Normalizing severities upon retrieval and insertion'''
        return value.casefold()

class NotificationConfig(WritableConfig):
    '''pydantic model representing the notification config'''
    _path = SNOOZE_CONFIG / 'notifications.yaml'

    notification_freq: int = 60
    notification_retry: int = 3

class HousekeeperConfig(WritableConfig):
    _path = SNOOZE_CONFIG / 'housekeeper.yaml'

    trigger_on_startup: bool = True
    record_ttl: timedelta = timedelta(days=2)
    cleanup_alert: timedelta = timedelta(minutes=5)
    cleanup_comment: timedelta = timedelta(days=1)
    cleanup_audit: timedelta = timedelta(days=28)
    cleanup_snooze: timedelta = timedelta(days=3)
    cleanup_notification: timedelta = timedelta(days=3)

class BackupConfig(ReadOnlyConfig):
    _path = SNOOZE_CONFIG / 'backup.yaml'

    enabled: bool = True
    path: Path = Path('/var/lib/snooze')
    excludes: List[str] = ('record', 'stats', 'comment', 'secrets')
