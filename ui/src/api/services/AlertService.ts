/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { AlertV2 } from '../models/AlertV2';
import type { Collection_AlertV2_ } from '../models/Collection_AlertV2_';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class AlertService {

  /**
   * Search
   * Return a list of alerts based on a search
   * @param pageNb
   * @param perPage
   * @param orderBy
   * @param direction
   * @param text
   * @returns Collection_AlertV2_ Successful Response
   * @throws ApiError
   */
  public static search(
    pageNb?: number,
    perPage: number = 20,
    orderBy?: string,
    direction?: number,
    text: string = '',
  ): CancelablePromise<Collection_AlertV2_> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/v2/alerts/search',
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
   * Get
   * Return a single alert by OID
   * @param oid
   * @returns AlertV2 Successful Response
   * @throws ApiError
   */
  public static get(
    oid: string,
  ): CancelablePromise<AlertV2> {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/v2/alerts/',
      query: {
        'oid': oid,
      },
      errors: {
        422: `Validation Error`,
      },
    });
  }

}
