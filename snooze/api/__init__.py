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
from pydantic import BaseModel, ValidationError

from snooze.api.routes import *
from snooze.utils.config import WebConfig
from snooze.health import HealthRoute
from snooze.token import TokenAuthMiddleware
from snooze.utils.typing import DuplicatePolicy, AuthorizationPolicy, RouteArgs, Pagination

log = getLogger('snooze.api')

# We're listing only the ones we might raise.
# It will not affect performance to list many, since falcon keep a dictionary of exception => handler.
USER_ERRORS = (
    # HTTP 400
    falcon.HTTPBadRequest, falcon.HTTPInvalidHeader, falcon.HTTPInvalidParam, falcon.HTTPMissingParam,
    # HTTP 401
    falcon.HTTPUnauthorized,
    # HTTP 403
    falcon.HTTPForbidden,
    # HTTP 404
    falcon.HTTPNotFound, falcon.HTTPRouteNotFound,
)

def log_warning_handler(err, _req, _resp, _params):
    '''Log caught exceptions as a warning'''
    log.warning(str(err), exc_info=err)
    raise err

def log_error_handler(err, _req, _resp, _params):
    '''Log caught exceptions as an error'''
    log.exception(err)
    raise err

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

        # Handler
        middlewares = [
            CORS(),
            LoggerMiddleware(self.core.config.core.audit_excluded_paths),
        ]
        if not self.core.config.core.no_login:
            middlewares += TokenAuthMiddleware(self.core.token_engine)
        self.handler = falcon.App(middleware=middlewares)
        self.handler.req_options.auto_parse_qs_csv = False

        json_handler = falcon.media.JSONHandler(
            dumps=bson.json_util.dumps,
            loads=bson.json_util.loads,
        )
        self.handler.req_options.media_handlers.update({'application/json': json_handler})
        self.handler.resp_options.media_handlers.update({'application/json': json_handler})
        self.handler.add_error_handler(USER_ERRORS, log_warning_handler)
        self.handler.add_error_handler(Exception, log_error_handler)
        self.auth_routes = {}
        # Alert route
        self.add_route('/alert', AlertRoute(self))
        # List route
        self.add_route('/login', LoginRoute(self))
        # Cluster route
        self.add_route('/cluster', ClusterRoute(self))
        # Plugin reload route
        self.add_route('/reload/{plugin_name}', ReloadPluginRoute(self))
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
        if self.core.config.ldap:
            self.auth_routes['ldap'] = LdapAuthRoute(self)
            self.add_route('/login/ldap', self.auth_routes['ldap'])
        # Optional metrics
        if self.core.stats.enabled:
            self.add_route('/metrics', MetricsRoute(self), '')

        web = self.core.config.core.web
        if web.enabled:
            prefix = '/web'
            self.add_route('/', RedirectRoute(), '')
            self.add_route(prefix, RedirectRoute(), '')
            sink_handler = StaticRoute(web.path, prefix).on_get
            self.handler.add_sink(sink_handler, prefix)

    def add_route(self, route, action, prefix='/api'):
        '''Map a falcon route class to a given path'''
        self.handler.add_route(prefix + route, action)

    def get_root_token(self):
        '''Return a root token for the root user. Used only when requesting it from the internal unix socket'''
        return self.jwt_auth.get_auth_token({'name': 'root', 'method': 'root', 'permissions': ['rw_all']})

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
