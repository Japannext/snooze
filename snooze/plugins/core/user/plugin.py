#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''The user plugin for managing user entries in the database'''

from datetime import datetime
from logging import getLogger

from snooze.plugins.core import Plugin
from snooze.utils.typing import AuthPayload, Profile, UserData

log = getLogger('snooze.user')

class User(Plugin):
    '''A plugin for managing user information in the database'''

    def update_user(self, auth: AuthPayload):
        '''Refresh a user's information in the database'''
        now = datetime.now().astimezone().strftime("%Y-%m-%dT%H:%M:%S%z")
        if auth.name == 'root' and auth.method == 'root':
            log.warning("Root user detected! Will not create an account in the database")
            return Profile(username=auth.username, method=auth.method)
        user_query = ['AND', ['=', 'name', auth.name], ['=', 'method', auth.method]]
        user_search = self.db.search('user', user_query)
        log.debug("Searching in users for user %s with method %s", auth.username, auth.method)
        if user_search['count'] > 0:
            user = user_search['data'][0]
            log.debug("User found: %s", user)
            old_groups = user.get('groups') or []
        else:
            log.debug("User not found, adding them to the database")
            user = auth.dict()
            old_groups = []
        user['last_login'] = now
        new_groups = auth.groups
        if old_groups != new_groups:
            log.debug("Will replace groups %s with %s", old_groups, new_groups)
            user['groups'] = new_groups
        if len(new_groups) > 0:
            query = ['IN', new_groups, 'groups']
            role_search = self.db.search('role', query)
            if role_search['count'] > 0:
                old_static_roles = user.get('static_roles') or []
                static_roles = list(map(lambda x: x['name'], role_search['data']))
                if old_static_roles != static_roles:
                    log.debug("Will replace static roles %s with %s", old_static_roles, static_roles)
                    user['static_roles'] = static_roles
                    user_roles = user.get('roles') or []
                    if user_roles:
                        log.debug("Will cleanup regular roles")
                        user['roles'] = [x for x in user_roles if x not in static_roles]

    def get_profile(self, auth: AuthPayload) -> Optional[Profile]:
        '''
        '''
        now = datetime.now().astimezone().strftime("%Y-%m-%dT%H:%M:%S%z")
        if auth.name == 'root' and auth.method == 'root':
            log.warning("Root user detected! Will not create an account in the database")
            return Profile(username=auth.username, method=auth.method)
        user_query = ['AND', ['=', 'name', auth.name], ['=', 'method', auth.method]]
        user_search = self.db.search('user', user_query)
        log.debug("Searching in users for user %s with method %s", name, method)
        if user_search['count'] > 0:
            user = user_search['data'][0]
            log.debug("User found: %s", user)
            old_groups = user.get('groups') or []
        else:
            log.debug("User not found, adding them to the database")
            user = auth_payload.dict()
            old_groups = []
        user['last_login'] = now
        new_groups = auth.groups
        if old_groups != new_groups:
            log.debug("Will replace groups %s with %s", old_groups, new_groups)
            user['groups'] = new_groups
        if len(new_groups) > 0:
            query = ['IN', new_groups, 'groups']
            role_search = self.db.search('role', query)
            if role_search['count'] > 0:
                old_static_roles = user.get('static_roles') or []
                static_roles = list(map(lambda x: x['name'], role_search['data']))
                if old_static_roles != static_roles:
                    log.debug("Will replace static roles %s with %s", old_static_roles, static_roles)
                    user['static_roles'] = static_roles
                    user_roles = user.get('roles') or []
                    if user_roles:
                        log.debug("Will cleanup regular roles")
                        user['roles'] = [x for x in user_roles if x not in static_roles]
        primary = self.meta.route_defaults.primary
        display_name = user.pop('display_name', '')
        email = user.pop('email', '')
        self.db.write('user', user, primary)
        profile_search = self.db.search('profile.general', user_query)
        if profile_search['count'] > 0:
            log.debug("User %s profile already exists, skipping", name)
            pref_search = self.db.search('profile.preferences', user_query)
            if pref_search['count'] > 0:
                return (profile_search['data'][0], pref_search['data'][0])
            else:
                return (profile_search['data'][0], None)
        else:
            log.debug("Creating user %s profile: Display Name (%s), Email (%s)", name, display_name, email)
            user_profile = {'name': name, 'method': method, 'display_name': display_name, 'email': email}
            self.db.write('profile.general', user_profile)
        return (None, None)
