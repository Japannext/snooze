#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Module for managing loading and writing the configuration files'''

import os
from abc import ABC
from contextlib import contextmanager
from logging import getLogger
from pathlib import Path
from datetime import timedelta
from typing import Optional, List, Any, Dict, Literal

import yaml
from filelock import FileLock
from pydantic import BaseModel, Field, PrivateAttr, validator

from snooze import __file__ as SNOOZE_PATH
from snooze.utils.typing import RouteArgs, HostPort

log = getLogger('snooze.utils.config')

SNOOZE_CONFIG = Path(os.environ.get('SNOOZE_SERVER_CONFIG', '/etc/snooze/server'))
SNOOZE_PLUGIN_PATH = Path(SNOOZE_PATH).parent / 'plugins/core'

class ReadOnlyConfig(BaseModel, ABC):
    '''A class representing a config file at a given path.
    Can only be read.'''
    _section: Optional[str] = PrivateAttr(None)
    _path: Optional[Path] = PrivateAttr(None)

    def __init__(self, basedir: Path = SNOOZE_CONFIG, data: Optional[dict] = None):
        section = self._class_get('_section')
        if section:
            self._class_set('_path', basedir / f"{section}.yaml")
        data = data or self._read()
        BaseModel.__init__(self, **data)

    def _class_get(self, key: str):
        '''Get a class attribute'''
        return getattr(self.__class__, key)

    def _class_set(self, key: str, value: Any):
        '''Set a class attribute'''
        # Using this workaround to avoid pydantic and WritableConfig own __setattr__
        setattr(self.__class__, key, value)

    def _read(self) -> dict:
        '''Read the config file and return the raw dict'''
        if self._path:
            try:
                return yaml.safe_load(self._path.read_text(encoding='utf-8'))
            except OSError:
                return {}
        else:
            return {}

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

    def __init__(self, basedir: Path = SNOOZE_CONFIG, data: Optional[dict] = None):
        ReadOnlyConfig.__init__(self, basedir, data)
        path = self._class_get('_path')
        path.touch(mode=0o600)
        self._class_set('_filelock', FileLock(path, timeout=1))

    @contextmanager
    def _lock(self):
        print(self)
        filelock = self._class_get('_filelock')
        filelock.acquire()
        self.refresh()
        try:
            yield # Update the config
        finally:
            self._update()
            filelock.release()

    def _update(self):
        '''Write a new config to the config file'''
        path = self._class_get('_path')
        if path:
            data = yaml.dump(self.dict())
            path.write_text(data, encoding='utf-8')

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
    class_name: str = Field('Route', alias='class')
    auto_reload: bool = False
    default_sorting: Optional[str] = None
    default_ordering: bool = True
    widgets: dict = Field(default_factory=dict)
    action_form: dict  = Field(default_factory=dict)
    audit: bool = True
    provides: List[str] = Field(default_factory=list)
    routes: Dict[str, RouteArgs] = Field(default_factory=dict)
    route_defaults: Optional[RouteArgs] = None
    icon: str = 'question-circle'
    options: dict = Field(default_factory=dict)
    search_fields: List[str] = Field(default_factory=list)

    def __init__(self, plugin_name: str):
        path = SNOOZE_PLUGIN_PATH / plugin_name / 'metadata.yaml'
        self._class_set('_path', path)
        ReadOnlyConfig.__init__(self)

class LdapConfig(WritableConfig):
    '''Configuration for LDAP authentication'''
    _section = 'ldap_auth'

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
    certfile: Optional[Path] = Field(None, env='SNOOZE_CERT_FILE')
    keyfile: Optional[Path] = Field(None, env='SNOOZE_KEY_FILE')

class WebConfig(BaseModel):
    '''The subconfig for the web server (snooze-web)'''
    enabled: bool = True
    path: Path = Path('/opt/snooze/web')

class CoreConfig(ReadOnlyConfig):
    '''Core configuration. Not editable live.'''
    _section = 'core'

    listen_addr: str = '0.0.0.0'
    port: int = 5200
    ssl: SslConfig = SslConfig()
    web: WebConfig = WebConfig()
    debug: bool = False
    bootstrap_db: bool = True
    unix_socket: Optional[Path] = Path('/var/run/snooze/server.socket')
    no_login: bool = False
    audit_excluded_paths: List[str] = ('/api/patlite', '/metrics', '/web')
    ssl: dict
    process_plugins: List[str] = ('rule', 'aggregaterule', 'snooze', 'notification')
    database: dict = Field(default_factory=lambda: {'type': 'file'})
    init_sleep: int = 5
    create_root_user: bool = False

class GeneralConfig(WritableConfig):
    '''General configuration of snooze'''
    _section = 'general'

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
    '''Configuration for default notification delays/retry'''
    _section = 'notifications'

    notification_freq: int = 60
    notification_retry: int = 3

class HousekeeperConfig(WritableConfig):
    '''Config for the housekeeper thread'''
    _section = 'housekeeper'

    trigger_on_startup: bool = True
    record_ttl: timedelta = timedelta(days=2)
    cleanup_alert: timedelta = timedelta(minutes=5)
    cleanup_comment: timedelta = timedelta(days=1)
    cleanup_audit: timedelta = timedelta(days=28)
    cleanup_snooze: timedelta = timedelta(days=3)
    cleanup_notification: timedelta = timedelta(days=3)

class BackupConfig(ReadOnlyConfig):
    '''Configuration for the backup job'''
    _section = 'backup'

    enabled: bool = True
    path: Path = Path('/var/lib/snooze')
    excludes: List[str] = ('record', 'stats', 'comment', 'secrets')

class ClusterConfig(ReadOnlyConfig):
    '''Configuration for the cluster'''
    _section = 'cluster'

    enabled: bool = False
    members: List[HostPort] = Field(tuple(HostPort(host='localhost')), env='SNOOZE_CLUSTER')

    @validator('members')
    def parse_members_env(cls, value): # pylint: disable=no-self-argument,no-self-use
        '''In case the environment (a string) is passed, parse the environment string'''
        if isinstance(value, str):
            members = []
            for member in value.split(','):
                members.append(HostPort(member.split(':', 1)))
            return members
        return value

