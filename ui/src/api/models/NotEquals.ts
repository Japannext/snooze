/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { AnyType } from './AnyType';
import type { DataField } from './DataField';

/**
 * Match if a field of a alert is not equal to a given value
 */
export type NotEquals = {
  kind?: '!=';
  field: DataField;
  value: AnyType;
};

