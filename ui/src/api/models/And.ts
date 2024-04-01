/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Condition } from './Condition';

/**
 * Match only if all conditions given in arguments match
 */
export type And = {
  kind?: 'and';
  conditions: Array<Condition>;
};

