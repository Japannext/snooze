/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Group } from './Group';

/**
 * A group of alerts, grouped together by a GroupingRule. Will be used
 * by alerts to reference which group they were grouped with
 */
export type AlertGroupV2 = {
  /**
   * Object ID (Mongodb)
   */
  oid?: string;
  /**
   * The group that all alerts identify with
   */
  group: Group;
  /**
   * The number of hits the alertgroup received during its lifetime
   */
  hits?: number;
  /**
   * Time of the last hit
   */
  last_hit?: string;
  /**
   * Message on the last alert received by the alert group
   */
  last_message?: string;
  /**
   * OID of the grouping rule that created this alert. If this value is set to None, the default group oid will be assumed.
   */
  grouping_rule_oid?: string;
};

