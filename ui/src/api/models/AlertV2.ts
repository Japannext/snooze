/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Group } from './Group';
import type { Severity } from './Severity';
import type { TraceContext } from './TraceContext';

/**
 * The description of an alert
 */
export type AlertV2 = {
  /**
   * Object ID (Mongodb)
   */
  oid?: string;
  timestamp?: string;
  observed_timestamp?: string;
  trace_context?: TraceContext;
  trace_flags?: string;
  severity?: Severity;
  /**
   * A map of key=value used to describe the resource which the alert comes from
   */
  resource?: any;
  /**
   * A map of key=value used to add contextual information to the alert
   */
  attributes?: any;
  /**
   * The body of the alert. Can be of multiple type
   */
  body: string;
  /**
   * Field related to snooze-process grouping feature. Is expected to be populated by snooze.process.grouping
   */
  group?: Group;
};

