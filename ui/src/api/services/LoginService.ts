/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { AuthSuccess } from '../models/AuthSuccess';
import type { BackendInfo } from '../models/BackendInfo';
import type { LdapCredentials } from '../models/LdapCredentials';
import type { TokenCredentials } from '../models/TokenCredentials';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class LoginService {

  /**
   * Info
   * List configured auth backends
   * @returns BackendInfo Successful Response
   * @throws ApiError
   */
  public static info(): CancelablePromise<Array<BackendInfo>> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/login/',
    });
  }

  /**
   * Auth Token
   * While not strictly necessary, this route helps making the
   * login action for token consistent with the other login methods
   * @param requestBody
   * @returns AuthSuccess Successful Response
   * @throws ApiError
   */
  public static authToken(
    requestBody: TokenCredentials,
  ): CancelablePromise<AuthSuccess> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/login/token',
      body: requestBody,
      mediaType: 'application/json',
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Auth Ldap
   * Authenticate with an OpenID Connect backend
   * @param name
   * @param requestBody
   * @returns any Successful Response
   * @throws ApiError
   */
  public static authLdap(
    name: string,
    requestBody: LdapCredentials,
  ): CancelablePromise<any> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/login/ldap/{name}',
      path: {
        'name': name,
      },
      body: requestBody,
      mediaType: 'application/json',
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Auth Oidc
   * Authenticate with an OpenID Connect backend
   * @param name
   * @returns any Successful Response
   * @throws ApiError
   */
  public static authOidc(
    name: string,
  ): CancelablePromise<any> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/login/oidc/{name}',
      path: {
        'name': name,
      },
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Auth Oauth2
   * Authenticate with an OAuth2 backend
   * @param name
   * @returns any Successful Response
   * @throws ApiError
   */
  public static authOauth2(
    name: string,
  ): CancelablePromise<any> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/login/oauth2/{name}',
      path: {
        'name': name,
      },
      errors: {
        422: `Validation Error`,
      },
    });
  }

}
