/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Condition } from './Condition';

/**
 * Match only if one of the condition given in arguments match
 */
export type Or = {
  kind?: 'or';
  conditions: Array<Condition>;
};

