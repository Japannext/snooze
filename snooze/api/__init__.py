#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''Module for handling the falcon WSGI'''

import os
import importlib.util
from abc import abstractmethod
from importlib import import_module
from logging import getLogger
from os.path import join as joindir
from pathlib import Path
from typing import Optional, List

import bson.json_util
import falcon
from falcon_auth import FalconAuthMiddleware, JWTAuthBackend
from falcon.errors import HTTPInternalServerError
from pydantic import BaseModel

from snooze.api.routes import *
from snooze.utils.config import WebConfig
from snooze.health import HealthRoute
from snooze.utils.typing import DuplicatePolicy, AuthorizationPolicy, RouteArgs, Pagination

log = getLogger('snooze.api')

class LoggerMiddleware:
    '''Middleware for logging'''

    def __init__(self, excluded_paths: List[str] = tuple()):
        self.logger = getLogger('snooze.audit')
        self.excluded_paths = excluded_paths

    def process_response(self, req, resp, *_args):
        '''Method for handling requests as a middleware'''
        source = req.access_route[0]
        method = req.method
        path = req.relative_uri
        status = resp.status[:3]
        message = f"{source} {method} {path} {status}"
        if not any(path.startswith(excluded) for excluded in self.excluded_paths):
            self.logger.debug(message)

class CORS:
    '''A falcon middleware to handle CORS when the snooze-server and
    snooze-web components are on different hosts.
    '''
    def __init__(self):
        pass

    def process_response(self, req, resp, _resource, req_succeeded):
        resp.set_header('Access-Control-Allow-Origin', '*')
        if (
            req_succeeded
            and req.method == 'OPTIONS'
            and req.get_header('Access-Control-Request-Method')
        ):
            allow = resp.get_header('Allow')
            resp.delete_header('Allow')

            allow_headers = req.get_header('Access-Control-Request-Headers', default='*')
            resp.set_headers(
                (
                    ('Access-Control-Allow-Methods', allow),
                    ('Access-Control-Allow-Headers', allow_headers),
                    ('Access-Control-Max-Age', '86400'),  # 24 hours
                )
            )

SNOOZE_GLOBAL_RUNDIR = '/var/run/snooze'
uid = os.getuid()
SNOOZE_LOCAL_RUNDIR = f"/var/run/user/{uid}"

class Api:
    def __init__(self, core: 'Core'):
        # Authentication
        self.core = core
        self.cluster = core.threads['cluster']

        # JWT setup
        if self.core.config.core.no_login:
            self.secret = ''
        else:
            self.secret = self.core.secrets['jwt_private_key']
        def auth(payload):
            log.debug("Payload received: %s", payload.get('user', {}).get('name', payload))
            return payload
        self.jwt_auth = JWTAuthBackend(auth, self.secret)

        # Handler
        middlewares = [
            CORS(),
            LoggerMiddleware(self.core.config.core.audit_excluded_paths),
            FalconAuthMiddleware(self.jwt_auth),
        ]
        self.handler = falcon.App(middleware=middlewares)
        self.handler.req_options.auto_parse_qs_csv = False

        json_handler = falcon.media.JSONHandler(
            dumps=bson.json_util.dumps,
            loads=bson.json_util.loads,
        )
        self.handler.req_options.media_handlers.update({'application/json': json_handler})
        self.handler.resp_options.media_handlers.update({'application/json': json_handler})
        self.handler.add_error_handler(Exception, self.custom_handle_uncaught_exception)
        self.auth_routes = {}
        # Alert route
        self.add_route('/alert', AlertRoute(self))
        # List route
        self.add_route('/login', LoginRoute(self))
        # Reload route
        self.add_route('/reload', ReloadRoute(self))
        # Cluster route
        self.add_route('/cluster', ClusterRoute(self))
        # Health route
        self.add_route('/health', HealthRoute(self))
        # Permissions route
        self.add_route('/permissions', PermissionsRoute(self))
        # Basic auth setup
        self.auth_routes['local'] = LocalAuthRoute(self)
        self.add_route('/login/local', self.auth_routes['local'])
        # Anonymous auth
        if self.core.config.general.anonymous_enabled:
            self.auth_routes['anonymous'] = AnonymousAuthRoute(self)
            self.add_route('/login/anonymous', self.auth_routes['anonymous'])
        # Ldap auth
        self.auth_routes['ldap'] = LdapAuthRoute(self)
        self.add_route('/login/ldap', self.auth_routes['ldap'])
        # Optional metrics
        if self.core.stats.enabled:
            self.add_route('/metrics', MetricsRoute(self), '')

        web = self.core.config.web
        if web.enabled:
            prefix = '/web'
            self.add_route('/', RedirectRoute(), '')
            self.add_route(prefix, RedirectRoute(), '')
            sink_handler = StaticRoute(web.path, prefix).on_get
            self.handler.add_sink(sink_handler, prefix)

    def custom_handle_uncaught_exception(self, e, req, resp, params):
        '''Custom handler for logging uncaught exceptions in falcon inside python logger.
        Make use of an internal method of falcon to do so.
        '''
        log.exception(e)
        self.handler._compose_error_response(req, resp, HTTPInternalServerError())

    def add_route(self, route, action, prefix='/api'):
        '''Map a falcon route class to a given path'''
        self.handler.add_route(prefix + route, action)

    def get_root_token(self):
        '''Return a root token for the root user. Used only when requesting it from the internal unix socket'''
        return self.jwt_auth.get_auth_token({'name': 'root', 'method': 'root', 'permissions': ['rw_all']})

    def reload(self, filename, auth_backends):
        '''Reload authentication backends configuration'''
        reloaded_auth = []
        reloaded_conf = []
        try:
            if self.core.reload_conf(filename):
                reloaded_conf.append(filename)
            for auth_backend in auth_backends:
                if self.auth_routes.get(auth_backend):
                    log.debug("Reloading %s auth backend", auth_backend)
                    self.auth_routes[auth_backend].reload()
                    reloaded_auth.append(auth_backend)
                else:
                    log.debug("Authentication backend '%s' not found", auth_backend)
            if len(reloaded_auth) > 0 or len(reloaded_conf) > 0:
                return {'status': falcon.HTTP_200, 'text': f"Reloaded auth '{reloaded_auth}' and conf {reloaded_conf}"}
            else:
                return {'status': falcon.HTTP_404, 'text': 'Error while reloading'}
        except Exception as e:
            log.exception(e)
            return {'status': falcon.HTTP_503}

    def write_and_reload(self, filename, conf, reload_conf, sync=False):
        '''Override the config files and reload. This is mainly used when changing the configuration
        from the web interface.
        '''
        result_dict = {}
        log.debug("Will write to %s config %s and reload %s", filename, conf, reload_conf)
        if filename and conf:
            res = write_config(filename, conf)
            if 'error' in res:
                return {'status': falcon.HTTP_503, 'text': res['error']}
            else:
                result_dict = {'status': falcon.HTTP_200, 'text': f"Reloaded config file {res['file']}"}
        if reload_conf:
            auth_backends = reload_conf.get('auth_backends', [])
            if auth_backends:
                result_dict = self.reload(filename, auth_backends)
            plugins = reload_conf.get('plugins', [])
            if plugins:
                result_dict = self.reload_plugins(plugins)
        if sync and self.cluster:
            self.cluster.write_and_reload(filename, conf, reload_conf)
        return result_dict

    def load_plugin_routes(self):
        log.debug('Loading plugin routes for API')
        for plugin in self.core.plugins:
            log.debug('Loading routes for %s at %s/%s/route.py', plugin.name, plugin.rootdir, self.api_type)
            spec = importlib.util.spec_from_file_location(
                f"snooze.plugins.core.{plugin.name}.{self.api_type}.route",
                joindir(plugin.rootdir, self.api_type, 'route.py')
            )
            plugin_module = importlib.util.module_from_spec(spec)
            try:
                spec.loader.exec_module(plugin_module)
                log.debug("Found custom routes for `%s`", plugin.name)
            except FileNotFoundError:
                # Loading default
                log.debug("Loading default route for `%s`", plugin.name)
                plugin_module = import_module(f"snooze.plugins.core.basic.{self.api_type}.route")
            except Exception as err:
                log.exception(err)
                log.debug("Skip loading plugin `%s` routes", plugin.name)
                continue
            for path, route_args in plugin.meta.routes:
                log.debug("For %s loading route: %s", path, route_args.dict())
                instance = getattr(plugin_module, route_args.class_name)(self, route_args)
                self.add_route(path, instance)

    def reload_plugins(self, plugins):
        '''Reload plugins'''
        plugins_error = []
        plugins_success = []
        log.debug("Reloading plugins %s", plugins)
        for plugin_name in plugins:
            plugin = self.core.get_core_plugin(plugin_name)
            if plugin:
                plugin.reload_data()
                plugins_success.append(plugin)
            else:
                plugins_error.append(plugin)
        if plugins_error:
            return {'status': falcon.HTTP_404, 'text': f"The following plugins could not be found: {plugin_error}"}
        else:
            return {'status': falcon.HTTP_200, 'text': "Reloaded plugins: {plugin_success}"}

