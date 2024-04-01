/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { ComparableType } from './ComparableType';
import type { DataField } from './DataField';

/**
 * Match if the field of an alert is strictly greater than a value.
 */
export type GreaterThan = {
  kind?: '>';
  field: DataField;
  value: ComparableType;
};

