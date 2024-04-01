/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

/**
 * An object allowing the alert to be grouped. This is usually
 * set by snooze-process during processing. When not set, the alert
 * won't be grouped
 */
export type Group = {
  /**
   * The hashed value of the aggregate field
   */
  hash: string;
  /**
   * A map of key=value used to group alerts together
   */
  attributes: any;
};

