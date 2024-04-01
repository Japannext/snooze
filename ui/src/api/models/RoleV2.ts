/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

/**
 * Represent a collection a permissions that can be assign to a user/user group
 */
export type RoleV2 = {
  /**
   * Object ID (Mongodb)
   */
  oid?: string;
  /**
   * Name of the role
   */
  name: string;
  /**
   * Human readable description of the role
   */
  description?: string;
  /**
   * A collection of scopes
   */
  scopes?: Array<string>;
};

