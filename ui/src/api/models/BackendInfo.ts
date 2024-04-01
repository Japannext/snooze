/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

/**
 * A structured data to return to the web interface in
 * order to display the correct form
 */
export type BackendInfo = {
  name: string;
  kind: BackendInfo.kind;
};

export namespace BackendInfo {

  export enum kind {
    LDAP = 'ldap',
    OIDC = 'oidc',
    OAUTH2 = 'oauth2',
  }


}

