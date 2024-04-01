/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Body_patch } from '../models/Body_patch';
import type { Collection_RuleV2_ } from '../models/Collection_RuleV2_';
import type { DocumentCreated } from '../models/DocumentCreated';
import type { RuleV2 } from '../models/RuleV2';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class RuleService {

  /**
   * Search
   * Return a list of rules based on a query
   * @param pageNb
   * @param perPage
   * @param orderBy
   * @param direction
   * @param text
   * @returns Collection_RuleV2_ Successful Response
   * @throws ApiError
   */
  public static search(
    pageNb?: number,
    perPage: number = 20,
    orderBy?: string,
    direction?: number,
    text: string = '',
  ): CancelablePromise<Collection_RuleV2_> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/v2/rule/search',
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
   * Validate
   * An endpoint to help with form validation
   * @param requestBody
   * @returns any Successful Response
   * @throws ApiError
   */
  public static validate(
    requestBody: RuleV2,
  ): CancelablePromise<any> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/v2/rule/validate',
      body: requestBody,
      mediaType: 'application/json',
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Get
   * Return only one rule object by OID
   * @param oid
   * @returns RuleV2 Successful Response
   * @throws ApiError
   */
  public static get(
    oid: string,
  ): CancelablePromise<RuleV2> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/v2/rule/{oid}',
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
   * Update one rule by OID
   * @param oid
   * @param requestBody
   * @returns any Successful Response
   * @throws ApiError
   */
  public static update(
    oid: string,
    requestBody: RuleV2,
  ): CancelablePromise<any> {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/v2/rule/{oid}',
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
   * Create a new rule
   * @param requestBody
   * @returns DocumentCreated Successful Response
   * @throws ApiError
   */
  public static create(
    requestBody: RuleV2,
  ): CancelablePromise<DocumentCreated> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/v2/rule/',
      body: requestBody,
      mediaType: 'application/json',
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Delete
   * Delete one or more rule by OID
   * @param requestBody
   * @returns any Successful Response
   * @throws ApiError
   */
  public static delete(
    requestBody: Array<string>,
  ): CancelablePromise<any> {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/v2/rule/',
      body: requestBody,
      mediaType: 'application/json',
      errors: {
        422: `Validation Error`,
      },
    });
  }

  /**
   * Patch
   * Patch multiple rules by OIDs
   * @param requestBody
   * @returns any Successful Response
   * @throws ApiError
   */
  public static patch(
    requestBody: Body_patch,
  ): CancelablePromise<any> {
    return __request(OpenAPI, {
      method: 'PATCH',
      url: '/v2/rule/',
      body: requestBody,
      mediaType: 'application/json',
      errors: {
        422: `Validation Error`,
      },
    });
  }

}
