/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Condition } from './Condition';

/**
 * Match the opposite of a given condition
 */
export type Not = {
  kind?: 'not';
  condition: Condition;
};

