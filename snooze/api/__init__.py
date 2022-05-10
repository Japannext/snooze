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
from pydantic import BaseModel, ValidationError

from snooze.api.routes import *
from snooze.utils.config import WebConfig
from snooze.health import HealthRoute
from snooze.token import TokenAuthMiddleware
from snooze.utils.typing import DuplicatePolicy, AuthorizationPolicy, RouteArgs, Pagination, USER_ERRORS
from snooze.utils.functions import log_error_handler, log_warning_handler, log_uncaught_handler

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

class Api(falcon.App):
    '''The WSGI class that handle all routes'''
    def __init__(self, core: 'Core'):
        # Authentication
        self.core = core
        self.cluster = core.threads['cluster']

        # Handler
        middlewares = [
            CORS(),
            LoggerMiddleware(self.core.config.core.audit_excluded_paths),
        ]
        if not self.core.config.core.no_login:
            middlewares.append(TokenAuthMiddleware(self.core.token_engine))
        falcon.App.__init__(self, middleware=middlewares)
        self.req_options.auto_parse_qs_csv = False

        json_handler = falcon.media.JSONHandler(
            dumps=bson.json_util.dumps,
            loads=bson.json_util.loads,
        )
        self.req_options.media_handlers.update({'application/json': json_handler})
        self.resp_options.media_handlers.update({'application/json': json_handler})
        self.add_error_handler(USER_ERRORS, log_warning_handler)
        self.add_error_handler(falcon.HTTPError, log_error_handler)
        self.add_error_handler(Exception, log_uncaught_handler)
        self.auth_routes = {}
        # Alert route
        self.add_route('/api/alert', AlertRoute(self))
        # List route
        self.add_route('/api/login', LoginRoute(self))
        # Cluster route
        self.add_route('/api/cluster', ClusterRoute(self))
        # Plugin reload route
        self.add_route('/api/reload/{plugin_name}', ReloadPluginRoute(self))
        # Health route
        self.add_route('/api/health', HealthRoute(self))
        # Permissions route
        self.add_route('/api/permissions', PermissionsRoute(self))
        # Basic auth setup
        self.auth_routes['local'] = LocalAuthRoute(self)
        self.add_route('/api/login/local', self.auth_routes['local'])
        # Anonymous auth
        if self.core.config.general.anonymous_enabled:
            self.auth_routes['anonymous'] = AnonymousAuthRoute(self)
            self.add_route('/api/login/anonymous', self.auth_routes['anonymous'])
        if self.core.config.ldap:
            self.auth_routes['ldap'] = LdapAuthRoute(self)
            self.add_route('/api/login/ldap', self.auth_routes['ldap'])
        # Optional metrics
        if self.core.stats.enabled:
            self.add_route('/metrics', MetricsRoute(self))

        web = self.core.config.core.web
        if web.enabled:
            prefix = '/web'
            self.add_route('/', RedirectRoute())
            self.add_route(prefix, RedirectRoute())
            sink_handler = StaticRoute(web.path, prefix).on_get
            self.add_sink(sink_handler, prefix)

    def get_root_token(self):
        '''Return a root token for the root user. Used only when requesting it from the internal unix socket'''
        auth = AuthPayload(username='root', method='root', permissions=['rw_all'])
        return self.core.token_engine.sign(auth)

    def load_plugin_routes(self):
        log.debug('Loading plugin routes for API')
        for plugin in self.core.plugins:
            log.debug("Loading routes for %s", plugin.name)
            spec = importlib.util.spec_from_file_location(
                f"snooze.plugins.core.{plugin.name}.falcon.route",
                joindir(plugin.rootdir, 'falcon', 'route.py')
            )
            plugin_module = importlib.util.module_from_spec(spec)
            try:
                spec.loader.exec_module(plugin_module)
                log.debug("Found custom routes for `%s`", plugin.name)
            except FileNotFoundError:
                # Loading default
                log.debug("Loading default route for `%s`", plugin.name)
                plugin_module = import_module('snooze.plugins.core.basic.falcon.route')
            except Exception as err:
                log.exception(err)
                log.debug("Skip loading plugin `%s` routes", plugin.name)
                continue
            log.debug('Routes: %s', plugin.meta.routes)
            for path, route_args in plugin.meta.routes.items():
                log.debug("For %s loading route: %s", path, route_args.dict())
                if route_args.class_name is not None:
                    instance = getattr(plugin_module, route_args.class_name)(self, plugin, route_args)
                    log.debug("Adding route %s: %s", path, instance)
                    self.add_route(path, instance)
