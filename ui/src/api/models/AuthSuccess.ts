/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

/**
 * A message to notify the user of the authentication success
 */
export type AuthSuccess = {
  /**
   * The token to use for all authenticated queries. Should be placed in the X-Snooze-Token header in HTTP requests
   */
  token: string;
};

