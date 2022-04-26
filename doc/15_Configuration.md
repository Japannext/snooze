# Summary

* [Core configuration](#core-configuration)
* [General configuration](#general-configuration)
* [Housekeeper configuration](#housekeeper-configuration)
* [Notification configuration](#notification-configuration)
* [LDAP configuration](#ldap-configuration)

# Core configuration

Core configuration. Not editable live. Require a restart of the server.
Usually located at `/etc/snooze/server/core.yaml`

## Properties

- `listen_addr` (string): IPv4 address on which Snooze process is listening to.  
**Default**: `0.0.0.0`.  
- `port` (integer): Port on which Snooze process is listening to.  
**Default**: `5200`.  
- `debug` (boolean): Activate debug log output.  
**Default**: `False`.  
- `bootstrap_db` (boolean): Populate the database with an initial configuration.  
**Default**: `True`.  
- `unix_socket` (string): Listen on this unix socket to issue root tokens.  
**Default**: `/var/run/snooze/server.socket`.  
- `no_login` (boolean): Disable Authentication (everyone has admin priviledges).  
**Environment variable**: `SNOOZE_NO_LOGIN`.  
**Default**: `False`.  
- `audit_excluded_paths` (array[string]): A list of HTTP paths excluded from audit logs. Any paththat starts with a path in this list will be excluded.  
**Default**: `['/api/patlite', '/metrics', '/web']`.  
- `process_plugins` (array[string]): List of plugins that will be used for processing alerts. Order matters.  
**Default**: `['rule', 'aggregaterule', 'snooze', 'notification']`.  
- `database` (object): Database.  
- `init_sleep` (integer): Time to sleep before retrying certain operations (bootstrap, clustering).  
**Default**: `5`.  
- `create_root_user` (boolean): Create a *root* user with a default password *root*.  
**Default**: `False`.  
- `ssl` ([SslConfig](#SslConfig)): SSL configuration.  
- `web` ([WebConfig](#WebConfig)): Web server configuration.  
- `cluster` ([ClusterConfig](#ClusterConfig)): Cluster configuration.  
- `backup` ([BackupConfig](#BackupConfig)): Backup configuration.  

## Definitions

### SslConfig

SSL configuration

#### Properties

- `enabled` (boolean): Enabling TLS termination.  
**Default**: `True`.  
- `certfile` (string): Path to the x509 PEM style certificate to use for TLS termination.  
**Environment variable**: `SNOOZE_CERT_FILE`.  
**Example(s)**:  
    - `/etc/pki/tls/certs/snooze.crt`
    - `/etc/ssl/certs/snooze.crt`

- `keyfile` (string): Path to the private key to use for TLS termination.  
**Environment variable**: `SNOOZE_KEY_FILE`.  
**Example(s)**:  
    - `/etc/pki/tls/private/snooze.key`
    - `/etc/ssl/private/snooze.key`


### WebConfig

The subconfig for the web server (snooze-web)

#### Properties

- `enabled` (boolean): Enable the web interface.  
**Default**: `True`.  
- `path` (string): Path to the web interface dist files.  
**Default**: `/opt/snooze/web`.  

### HostPort

An object to represent a host-port pair

#### Properties

- `host` (string) (**required**): The host address to reach (IP or resolvable hostname).  
- `port` (integer): The port where the host is expected to listen to.  
**Default**: `5200`.  

### ClusterConfig

Configuration for the cluster

#### Properties

- `enabled` (boolean): Enable clustering. Required when running multiple backends.  
**Default**: `False`.  
- `members` (array): List of snooze servers in the cluster. If the environment variable is provided, a special syntax is expected (`"<host>:<port>,<host>:<port>,..."`).  
**Environment variable**: `SNOOZE_CLUSTER`.  
**Example(s)**:  
    - `[{'host': 'host01', 'port': 5200}, {'host': 'host02', 'port': 5200}, {'host': 'host03', 'port': 5200}]`
    - `host01:5200,host02:5200,host03:5200`


### BackupConfig

Configuration for the backup job

#### Properties

- `enabled` (boolean): Enable backups.  
**Default**: `True`.  
- `path` (string): Path to store database backups.  
**Default**: `/var/lib/snooze`.  
- `excludes` (array[string]): Collections to exclude from backups.  
**Default**: `['record', 'stats', 'comment', 'secrets']`.  


# General configuration

General configuration of snooze. Can be edited live in the web interface.
Usually located at `/etc/snooze/server/general.yaml`.

## Properties

- `default_auth_backend` (string): Backend that will be first in the list of displayed authentication backends.  
**Default**: `local`.  
- `local_users_enabled` (boolean): Enable the creation of local users in snooze. This can be disabled when another reliable authentication backend is used, and the admin want to make auditing easier.  
**Default**: `True`.  
- `metrics_enabled` (boolean): Enable Prometheus metrics.  
**Default**: `True`.  
- `anonymous_enabled` (boolean): Enable anonymous user login. When a user log in as anonymous, he will be given user permissions.  
**Default**: `False`.  
- `ok_severities` (array[string]): List of severities that will automatically close the aggregate upon entering the system. This is mainly for icinga/grafana that can close the alert when the status becomes green again.  
**Default**: `['ok', 'success']`.  

# Housekeeper configuration

Config for the housekeeper thread. Can be edited live in the web interface.
Usually located at `/etc/snooze/server/housekeeper.yaml`.

## Properties

- `trigger_on_startup` (boolean): Trigger all housekeeping job on startup.  
**Default**: `True`.  
- `record_ttl` (number): Default TTL (in seconds) for alerts incoming.  
**Default**: `172800.0`.  
- `cleanup_alert` (number): Time (in seconds) between each run of alert cleaning. Alerts that exceeded their TTL  will be deleted.  
**Default**: `300.0`.  
- `cleanup_comment` (number): Time (in seconds) between each run of comment cleaning. Comments which are not bound to any alert will be deleted.  
**Default**: `86400.0`.  
- `cleanup_audit` (number): Cleanup orphans audit logs that are older than the given duration (in seconds). Run daily.  
**Default**: `2419200.0`.  
- `cleanup_snooze` (number): Cleanup snooze filters that have been expired for the given duration (in seconds). Run daily.  
**Default**: `259200.0`.  
- `cleanup_notification` (number): Cleanup notifications that have been expired for the given duration (in seconds). Run daily.  
**Default**: `259200.0`.  

# Notification configuration

Configuration for default notification delays/retry. Can be edited live in the web interface.
Usually located at `/etc/snooze/server/notifications.yaml`.

## Properties

- `notification_freq` (number): Time (in seconds) to wait before sending the next notification.  
**Default**: `60.0`.  
- `notification_retry` (integer): Number of times to retry sending a failed notification.  
**Default**: `3`.  

# LDAP configuration

Configuration for LDAP authentication. Can be edited live in the web interface.
Usually located at `/etc/snooze/server/ldap_auth.yaml`.

## Properties

- `enabled` (boolean): Enable or disable LDAP Authentication.  
**Default**: `False`.  
- `base_dn` (string) (**required**): LDAP users location. Multiple DNs can be added if separated by semicolons.  
- `user_filter` (string) (**required**): LDAP search filter for the base DN.  
**Example(s)**:  
    - `(objectClass=posixAccount)`

- `bind_dn` (string) (**required**): Distinguished name to bind to the LDAP server.  
**Example(s)**:  
    - `CN=john.doe,OU=users,DC=example,DC=com`

- `bind_password` (string) (**required**): Password for the Bind DN user.  
- `host` (string) (**required**): LDAP host.  
**Example(s)**:  
    - `ldaps://example.com`

- `port` (integer): LDAP server port.  
**Default**: `636`.  
- `group_dn` (string): Base DN used to filter out groups. Will default to the User base DN Multiple DNs can be added if separated by semicolons.  
- `email_attribute` (string): User attribute that displays the user email adress.  
**Default**: `mail`.  
- `display_name_attribute` (string): User attribute that displays the user real name.  
**Default**: `cn`.  
- `member_attribute` (string): Member attribute that displays groups membership.  
**Default**: `memberof`.  

