/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Collection_RoleV2_ } from '../models/Collection_RoleV2_';
import type { DocumentCreated } from '../models/DocumentCreated';
import type { RoleV2 } from '../models/RoleV2';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class RoleService {

  /**
   * Search
   * Return a list of roles based on a query
   * @param pageNb
   * @param perPage
   * @param orderBy
   * @param direction
   * @param text
   * @returns Collection_RoleV2_ Successful Response
   * @throws ApiError
   */
  public static search(
    pageNb?: number,
    perPage: number = 20,
    orderBy?: string,
    direction?: number,
    text: string = '',
  ): CancelablePromise<Collection_RoleV2_> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/v2/role/search',
      query: {
        'page_nb': pageNb,
        'per_page': perPage,
        'order_by': orderBy,
        'direction': direction,
        'text': text,
      },
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Permissions
   * Return the list of possible permissions
   * @returns string Successful Response
   * @throws ApiError
   */
  public static permissions(): CancelablePromise<Array<string>> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/v2/role/permissions',
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Validate
   * An endpoint to help with form validation
   * @param requestBody
   * @returns any Successful Response
   * @throws ApiError
   */
  public static validate(
    requestBody: RoleV2,
  ): CancelablePromise<any> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/v2/role/validate',
      body: requestBody,
      mediaType: 'application/json',
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Get
   * Return only one role object by OID
   * @param oid
   * @returns RoleV2 Successful Response
   * @throws ApiError
   */
  public static get(
    oid: string,
  ): CancelablePromise<RoleV2> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/v2/role/{oid}',
      path: {
        'oid': oid,
      },
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Update
   * Update one role by OID
   * @param oid
   * @param requestBody
   * @returns any Successful Response
   * @throws ApiError
   */
  public static update(
    oid: string,
    requestBody: RoleV2,
  ): CancelablePromise<any> {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/v2/role/{oid}',
      path: {
        'oid': oid,
      },
      body: requestBody,
      mediaType: 'application/json',
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Create
   * Create a new role
   * @param requestBody
   * @returns DocumentCreated Successful Response
   * @throws ApiError
   */
  public static create(
    requestBody: RoleV2,
  ): CancelablePromise<DocumentCreated> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/v2/role/',
      body: requestBody,
      mediaType: 'application/json',
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Delete
   * Delete one or more role by OID
   * @param requestBody
   * @returns any Successful Response
   * @throws ApiError
   */
  public static delete(
    requestBody: Array<string>,
  ): CancelablePromise<any> {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/v2/role/',
      body: requestBody,
      mediaType: 'application/json',
      errors: {
        422: `Validation Error`,
      },
    });
  }

}
