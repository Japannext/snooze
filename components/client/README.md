# Snooze client

CLI and python API to [Snooze server](https://github.com/snoozeweb/snooze).


## Installation

```bash
pip3 install --user snooze-client
```

## Configuration

Snooze client will load the configuration at `/etc/snooze/client.yaml`.
You can override the configuration path by using the `SNOOZE_CLIENT_CONFIG_FILE` environment variable.

Mandatory configuration:
* `server` (String): URL to the Snooze server. Example: `https://snooze.example.com:5200`

Authentication config (required for endpoints that require authorization):
* `auth_method` (Possible values: `local`, `ldap`): The authentication backup to use to connect to the snooze server. The `local` backend is for local user, the `ldap` backend is for connecting with a LDAP user.
* `credentials`: The data required by an authentication backend to connect.
  * For `local`/`ldap`:
    * `username` (String): The name of the user
    * `password` (String): The password of the user

Optional config:
* `ca_bundle` (String): Path to a custom CA bundle to use for TLS verification. The default will use python's CAs and the OS CAs.
* `app_name` (String): Name of the application to appear as when doing actions through the API. Defaults to `snooze_client`.
* `token_to_disk` (Boolean): Set to `false` if you do not want to cache the token to disk between calls. Defaults to `false`. The token is written to one of the following location (in order, if directory exists): the `SNOOZE_TOKEN_PATH` environment variable, `${PWD}/.snooze-token`, `${HOME}/.snooze-token`.

## Usage

### Add a snooze entry

CLI:
```bash
snooze_client snooze -q "process=myapp and hostname=${HOSTNAME}" -t "$(date -Is)" "$(date -Is -d '1 hour')"
```

HTTP API equivalent:
```bash
# Password will be fetched from prompt
TOKEN="$(curl -u $USER https://snooze.example.com/api/login/local | jq .token -r)"
cat <<payload EOF
{
    "name": "[snooze_client] snooze entry",
    "qls": [{"ql": "", "field": "condition"}],
    "time_constraint": {
        "weekdays": ["Monday", 2, "Wed"],
        "time": [
            {"from": "09:00", "until": "10:00"},
        ],
        "datetime": [
            {"from": "2021-07-01T12:00:00", "until": "2021-08-31T12:00:00"},
        ],
    },
    "comment": "Created by curl",
}
EOF
curl \
    -X POST
    -H "Authorization: JWT ${TOKEN}" \
    -H "Content-Type: application/json" \
    https://snooze.example.com/api/snooze \
    -d @payload
```
