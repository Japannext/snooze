/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

/**
 * The body to send inside a LDAP login request
 */
export type LdapCredentials = {
  /**
   * The username of the user
   */
  username: string;
  /**
   * The password of the user
   */
  password: string;
};

