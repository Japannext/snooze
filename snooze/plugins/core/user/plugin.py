#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''The user plugin for managing user entries in the database'''

from datetime import datetime
from logging import getLogger
from typing import Tuple

from snooze.plugins.core.basic.plugin import Plugin
from snooze.utils.typing import AuthPayload, ProfileGeneral, ProfilePreferences

log = getLogger('snooze.user')

class UserPlugin(Plugin):
    '''A plugin for managing user information in the database'''

    def update_user(self, auth: AuthPayload):
        '''Refresh a user's information in the database'''
        now = datetime.now().astimezone().strftime("%Y-%m-%dT%H:%M:%S%z")
        if auth.username == 'root' and auth.method == 'root':
            log.warning("Root user detected! Will not create an account in the database")
            return
        user_query = ['AND', ['=', 'name', auth.username], ['=', 'method', auth.method]]
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

    def get_profile(self, auth: AuthPayload) -> Tuple[ProfileGeneral, ProfilePreferences]:
        '''Return the profile options for the user'''
        if auth.username == 'root' and auth.method == 'root':
            return (
                ProfileGeneral(username=auth.username, method=auth.method),
                ProfilePreferences(username=auth.username, method=auth.method),
            )

        user_query = ['AND', ['=', 'name', auth.username], ['=', 'method', auth.method]]
        profiles_general = self.db.search('profile.general', user_query)['data']
        profiles_preferences = self.db.search('profile.preferences', user_query)['data']
        if profiles_general:
            general = profiles_general[0]
        else:
            general = ProfileGeneral(name=auth.username, method=auth.method, display_name=auth.username)
            self.db.write('profile.general', general.dict())
        if profiles_preferences:
            preferences = profiles_preferences[0]
        else:
            preferences = ProfilePreferences(name=auth.username, method=auth.method)
        return general, preferences
