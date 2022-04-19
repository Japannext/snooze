#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Module for LDAP related functionalities'''

from logging import getLogger

import falcon
from ldap3 import Server, Connection, ALL, SUBTREE
from ldap3.core.exceptions import LDAPOperationResult, LDAPExceptionError

from snooze.api.falcon import AuthRoute
from snooze.utils.config import LdapConfig

log = getLogger('snooze.api.ldap')

class LdapAuthRoute(AuthRoute):
    '''An authentication route for LDAP users'''
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.name = 'Ldap'
        self.enabled = False
        self.config = LdapConfig()
        self.reload()

    def reload(self):
        self.config.refresh()
        self.enabled = self.config.enabled
        if self.enabled:
            try:
                if '://' in self.config.host:
                    uri = self.config.host
                else:
                    if self.config.port == 636:
                        uri = f"ldaps://{self.config.host}"
                    else:
                        uri = f"ldap://{self.config.host}"
                self.server = Server(uri, port=self.config.port, get_info=ALL, connect_timeout=10)
                bind_con = Connection(
                    self.server,
                    user=self.config.bind_dn,
                    password=self.config.bind_password,
                    raise_exceptions=True
                )
                if not bind_con.bind():
                    log.error("Cannot BIND to LDAP server: %s:%s", uri, self.config.port)
                    self.enabled = False
            except Exception as err:
                log.exception(err)
                self.enabled = False
        log.debug("Authentication backend 'ldap'. Enabled: %s", self.config.enabled)

    def _search_user(self, username):
        try:
            bind_con = Connection(
                self.server,
                user=self.config.bind_dn,
                password=self.config.bind_password,
                raise_exceptions=True
            )
            bind_con.bind()
            user_filter = self.config.user_filter.replace('%s', username)
            bind_con.search(
                search_base = self.config.base_dn,
                search_filter = user_filter,
                attributes = [
                    self.config.display_name_attribute,
                    self.config.email_attribute,
                    self.config.member_attribute,
                ],
                search_scope = SUBTREE
            )
            response = bind_con.response
            if (
                bind_con.result['result'] == 0
                and len(response) > 0
                and 'dn' in response[0].keys()
            ):
                user_dn = response[0]['dn']
                attributes = response[0]['attributes']
                groups = [
                    group for group in attributes[self.config.member_attribute]
                    for dn in self.config.group_dn.split(':')
                    if group.endswith(dn)
                ]
                return {'name': username, 'dn': user_dn, 'groups': groups}
            else:
                # Could not find user in search
                raise falcon.HTTPUnauthorized(description="Error in search: Could not find" \
                    f" user {username} in LDAP search")
        except LDAPOperationResult as err:
            raise falcon.HTTPUnauthorized(description=f"Error during search: {err}")
        except LDAPExceptionError as err:
            raise falcon.HTTPUnauthorized(description=f"Error during search: {err}")

    def _bind_user(self, user_dn, password):
        try:
            user_con = Connection(
                self.server,
                user=user_dn,
                password=password,
                raise_exceptions=True
            )
            user_con.bind()
            return user_con
        except LDAPOperationResult as err:
            raise falcon.HTTPUnauthorized(description=f"Error during bind: {err}")
        except LDAPExceptionError as err:
            raise falcon.HTTPUnauthorized(description=f"Error during bind: {err}")
        finally:
            user_con.unbind()

    def authenticate(self, req, resp):
        username, password = self._extract_credentials(req)
        user = self._search_user(username)
        user_con = self._bind_user(user['dn'], password)
        if user_con.result['result'] == 0:
            req.context['user'] = user
        else:
            raise falcon.HTTPUnauthorized(description="")

    def parse_user(self, user):
        groups = list(map(lambda x: x.split(',')[0].split('=')[1], user['groups']))
        return {'name': user['name'], 'groups': groups, 'method': 'ldap'}
