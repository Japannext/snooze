/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { ComparableType } from './ComparableType';
import type { DataField } from './DataField';

/**
 * Match if the field of an alert is greater than or equal to a value.
 */
export type LowerOrEqual = {
  kind?: '<=';
  field: DataField;
  value: ComparableType;
};

