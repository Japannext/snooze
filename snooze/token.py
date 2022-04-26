#
# Copyright 2018-2020 Florian Dematraz <florian.dematraz@snoozeweb.net>
# Copyright 2018-2020 Guillaume Ludinard <guillaume.ludi@gmail.com>
# Copyright 2020-2021 Japannext Co., Ltd. <https://www.japannext.co.jp/>
# SPDX-License-Identifier: AFL-3.0
#

'''A module for managing the token engine'''

import falcon
import jwt
from jwt.exceptions import InvalidTokenError

class TokenEngine:
    '''Sign and verify tokens'''
    def __init__(self, secret_key, algorithm='HS256'):
        self.secret = secret_key
        self.algorithm = algorithm

    def sign(self, payload):
        '''Sign a payload and return the token'''
        token = jwt.encode(payload, self.secret, algorithm=self.algorithm)
        return token

    def verify(self, token):
        '''Verify the token and return the payload'''
        payload = jwt.decode(token, self.secret, algorithm=[self.algorithm])
        return payload

class TokenAuthMiddleware:
    '''A falcon middleware for '''
    def __init__(self, engine: TokenEngine):
        self.scheme = 'JWT'
        self.engine = engine
    def process_resource(self, req, resp, resource, *args, **kwargs):
        authorization = req.get_header('Authorization')
        try:
            scheme, credentials = authorization.split(' ', 1)
        except ValueError as err:
            raise falcon.HTTPInvalidHeader(header_name='Authorization',
                description=f"Must be in the form `{self.scheme} <credentials>`") from err
        if scheme != self.scheme:
            raise falcon.HTTPUnauthorized(description=f"Invalid authorization scheme: {scheme}."
                f"Must be {self.scheme}")
        try:
            payload = self.engine.verify(credentials)
        except InvalidTokenError as err:
            raise falcon.HTTPUnauthorized(header_name='Authorization',
                message=str(err)) from err

        req.context['auth'] = AuthPayload(**payload)
