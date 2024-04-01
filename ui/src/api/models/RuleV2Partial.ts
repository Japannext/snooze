/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Condition } from './Condition';

/**
 * A rule is an object describin how to transform an incoming
 * alert if it matches a condition
 */
export type RuleV2Partial = {
  /**
   * Object ID (Mongodb)
   */
  oid?: string;
  /**
   * Name of the rule
   */
  name?: string | null;
  /**
   * Human readable description of the rule
   */
  description?: string | null;
  /**
   * A condition in which the transformation will be triggered
   */
  condition?: Condition | null;
};

