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
from pydantic import BaseModel, Field, PrivateAttr, validator, ValidationError
from pydantic.dataclasses import dataclass

from snooze import __file__ as SNOOZE_PATH
from snooze.utils.typing import RouteArgs, HostPort, Widget

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
        data = data or self._read() or {}
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
        if data is None:
            data = {}
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
            object.__setattr__(self, key, value)

    def update(self, values: dict):
        '''Update the config with a dictionary'''
        with self._lock():
            for key, value in values.items():
                object.__setattr__(self, key, value)

class MetadataConfig(ReadOnlyConfig):
    '''A class to fetch metadata configuration'''
    name: Optional[str] = None
    desc: Optional[str] = None
    class_name: Optional[str] = Field('Route', alias='class')
    auto_reload: bool = False
    default_sorting: Optional[str] = None
    default_ordering: bool = True
    audit: bool = True
    widgets: Dict[str, Widget] = Field(default_factory=dict)
    action_form: dict  = Field(default_factory=dict)
    provides: List[str] = Field(default_factory=list)
    routes: Dict[str, RouteArgs] = Field(default_factory=dict)
    route_defaults: RouteArgs = RouteArgs()
    icon: str = 'question-circle'
    options: dict = Field(default_factory=dict)
    search_fields: List[str] = Field(default_factory=list)

    def __init__(self, plugin_name: str):
        path = SNOOZE_PLUGIN_PATH / plugin_name / 'metadata.yaml'
        self._class_set('_path', path)
        data = self._read() or {}
        try:
            BaseModel.__init__(self, **data)
        except ValidationError as err:
            raise Exception(f"Cannot load metadata for plugin {plugin_name}") from err

class LdapConfig(WritableConfig):
    '''Configuration for LDAP authentication'''
    _section = 'ldap_auth'

    enabled: bool = Field(
        description='Enable or disable LDAP Authentication',
        default=False,
    )
    base_dn: str = Field(
        title='Base DN',
        description='LDAP users location. Multiple DNs can be added if separated by semicolons',
    )
    user_filter: str = Field(
        title='User base filter',
        description='LDAP search filter for the base DN (eg. (objectClass=posixAccount))',
    )
    bind_dn: str = Field(
        title='Bind DN',
        description='Distinguished name to bind to the LDAP server (eg. '
        'CN=john.doe,OU=users,DC=example,DC=com)',
    )
    bind_password: str = Field(
        title='Bind DN password',
        description='Password for the Bind DN user',
        exclude=True,
    )
    host: str = Field(
        description='LDAP host (eg. ldaps://example.com)',
    )
    port: int = Field(
        default=636,
        description='LDAP server port (389 by default if no SSL, 636 if SSL is enabled)',
    )
    group_dn: Optional[str] = Field(
        title='Group DN',
        default=None,
        description='Base DN used to filter out groups. Will default to the User base DN'
        ' Multiple DNs can be added if separated by semicolons',
    )
    email_attribute: str = Field(
        title='Email attribute',
        default='mail',
        description='User attribute that displays the user email adress',
    )
    display_name_attribute: str = Field(
        title='Display name attribute',
        default='cn',
        description='User attribute that displays the user real name',
    )
    member_attribute: str = Field(
        title='Member attribute',
        default='memberof',
        description='Member attribute that displays groups membership',
    )

class SslConfig(BaseModel):
    '''SSL configuration'''
    enabled: bool = Field(
        default=True,
        description='Enabling TLS termination',
    )
    certfile: Optional[Path] = Field(
        title='Certificate file',
        env='SNOOZE_CERT_FILE',
        default=None,
        description='Path to the x509 PEM style certificate to use for TLS termination',
    )
    keyfile: Optional[Path] = Field(
        title='Key file',
        default=None,
        env='SNOOZE_KEY_FILE',
        description='Path to the private key to use for TLS termination',
    )

class WebConfig(BaseModel):
    '''The subconfig for the web server (snooze-web)'''
    enabled: bool = Field(
        default=True,
        description='Enable the web interface',
    )
    path: Path = Field(
        default='/opt/snooze/web',
        description='Path to the web interface dist files',
    )

class CoreConfig(ReadOnlyConfig):
    '''Core configuration. Not editable live.'''
    _section = 'core'

    listen_addr: str = Field(
        title='Listening address',
        default='0.0.0.0',
        description="IPv4 address on which Snooze process is listening to",
    )
    port: int = Field(
        default=5200,
        description='Port on which Snooze process is listening to',
    )
    debug: bool = Field(
        default=False,
        description='Activate debug log output',
    )
    bootstrap_db: bool = Field(
        title='Bootstrap database',
        default=True,
        description='Populate the database with an initial configuration',
    )
    unix_socket: Optional[Path] = Field(
        title='Unix socket',
        default='/var/run/snooze/server.socket',
        description='Listen on this unix socket to issue root tokens',
    )
    no_login: bool = Field(
        title='No login',
        default=False,
        env='SNOOZE_NO_LOGIN',
        description='Disable Authentication (everyone has admin priviledges)',
    )
    audit_excluded_paths: List[str] = Field(
        title='Audit excluded paths',
        default=('/api/patlite', '/metrics', '/web'),
        description='A list of HTTP paths excluded from audit logs. Any path'
        'that starts with a path in this list will be excluded.',
    )
    process_plugins: List[str] = Field(
        title='Process plugins',
        default=('rule', 'aggregaterule', 'snooze', 'notification'),
        description='List of plugins that will be used for processing alerts.'
        ' Order matters.',
    )
    database: dict = Field(
        title='Database',
        default_factory=lambda: {'type': 'file'},
    )
    init_sleep: int = Field(
        title='Init sleep',
        default=5,
        description='Time to sleep before retrying certain operations (bootstrap, clustering)',
    )
    create_root_user: bool = Field(
        title='Create root user',
        default=False,
        description='Create a *root* user with a default password *root*',
    )
    ssl: SslConfig = Field(title='SSL configuration', default_factory=SslConfig)
    web: WebConfig = Field(title='Web server configuration', default_factory=WebConfig)

class GeneralConfig(WritableConfig):
    '''General configuration of snooze'''
    _section = 'general'

    default_auth_backend: Literal['local', 'ldap'] = Field(
        title='Default authentication backend',
        description='Backend that will be first in the list of displayed authentication backends',
        default='local',
    )
    local_users_enabled: bool = Field(
        title='Local users enabled',
        description='Enable the creation of local users in snooze. This can be disabled when another'
        ' reliable authentication backend is used, and the admin want to make auditing easier',
        default=True,
    )
    metrics_enabled: bool = Field(
        title='Metrics enabled',
        description='Enable Prometheus metrics',
        default=True,
    )
    anonymous_enabled: bool = Field(
        title='Anonymous enabled',
        description='Enable anonymous user login. When a user log in as anonymous, he will be given user permissions',
        default=False,
    )
    ok_severities: List[str] = Field(
        title='OK severities',
        description='List of severities that will automatically close the aggregate upon entering the system.'
        ' This is mainly for icinga/grafana that can close the alert when the status becomes green again',
        default=('ok', 'success'),
    )

    @validator('ok_severities', each_item=True)
    def normalize_severities(cls, value): # pylint: disable=no-self-argument,no-self-use
        '''Normalizing severities upon retrieval and insertion'''
        return value.casefold()

class NotificationConfig(WritableConfig):
    '''Configuration for default notification delays/retry'''
    _section = 'notifications'

    notification_freq: int = Field(
        title='Frequency',
        description='Time to wait before sending the next notification',
        default=60,
    )
    notification_retry: int = Field(
        title='Retry number',
        description='',
        default=3,
    )

class HousekeeperConfig(WritableConfig):
    '''Config for the housekeeper thread'''
    _section = 'housekeeper'

    trigger_on_startup: bool = Field(
        title='Trigger on startup',
        default=True,
        description='Trigger all housekeeping job on startup',
    )
    record_ttl: timedelta = Field(
        title='Record Time-To-Live',
        description='Default TTL for alerts incoming',
        default=timedelta(days=2),
    )
    cleanup_alert: timedelta = Field(
        title='Cleanup alert',
        description='Time between each run of alert cleaning. Alerts that exceeded their TTL will be deleted',
        default=timedelta(minutes=5),
    )
    cleanup_comment: timedelta = Field(
        title='Cleanup comment',
        description='Time between each run of comment cleaning. Comments which are not bound to any alert will'
        ' be deleted',
        default=timedelta(days=1),
    )
    cleanup_audit: timedelta = Field(
        title='Cleanup audit',
        description='Cleanup orphans audit logs that are older than the given duration. Run daily',
        default=timedelta(days=28),
    )
    cleanup_snooze: timedelta = Field(
        title='Cleanup snooze',
        description="Cleanup snooze filters that have been expired for the given duration. Run daily",
        default=timedelta(days=3),
    )
    cleanup_notification: timedelta = Field(
        title='Cleanup notifications',
        description='Cleanup notifications that have been expired for the given duration. Run daily',
        default=timedelta(days=3),
    )

    class Config:
        json_encoders = {
            # timedelta should be serialized into seconds (int)
            timedelta: lambda dt: dt.total_seconds(),
        }

class BackupConfig(ReadOnlyConfig):
    '''Configuration for the backup job'''
    _section = 'backup'

    enabled: bool = Field(
        default=True,
        description='Enable backups',
    )
    path: Path = Field(
        default='/var/lib/snooze',
        description='Path to store database backups',
    )
    excludes: List[str] = Field(
        description='Collections to exclude from backups',
        default=('record', 'stats', 'comment', 'secrets'),
    )

class ClusterConfig(ReadOnlyConfig):
    '''Configuration for the cluster'''
    _section = 'cluster'

    enabled: bool = Field(
        default=False,
        description='Enable clustering. Required when running multiple backends',
    )
    members: List[HostPort] = Field(
        env='SNOOZE_CLUSTER',
        default=tuple(HostPort(host='localhost')),
        description='List of snooze servers in the cluster {host, port}'
    )

    @validator('members')
    def parse_members_env(cls, value): # pylint: disable=no-self-argument,no-self-use
        '''In case the environment (a string) is passed, parse the environment string'''
        if isinstance(value, str):
            members = []
            for member in value.split(','):
                members.append(HostPort(member.split(':', 1)))
            return members
        return value

@dataclass
class Config:
    '''An object representing the complete snooze configuration'''
    basedir: Path

    core: CoreConfig
    general: GeneralConfig
    housekeeper: HousekeeperConfig
    cluster: ClusterConfig
    backup: BackupConfig
    ldap: Optional[LdapConfig]

    def __init__(self, basedir: Path = SNOOZE_CONFIG):
        self.basedir = basedir
        self.core = CoreConfig(basedir)
        self.general = GeneralConfig(basedir)
        self.notifications = NotificationConfig(basedir)
        self.housekeeper = HousekeeperConfig(basedir)
        self.cluster = ClusterConfig(basedir)
        self.backup = BackupConfig(basedir)
        try:
            self.ldap = LdapConfig(basedir)
        except (FileNotFoundError, ValidationError):
            self.ldap = None
