/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { ComparableType } from './ComparableType';
import type { DataField } from './DataField';

/**
 * Match if the field of an alert is strictly lower than a value.
 */
export type LowerThan = {
  kind?: '<';
  field: DataField;
  value: ComparableType;
};

