/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { AnyType } from './AnyType';
import type { DataField } from './DataField';

/**
 * Match if the field of a alert is exactly equal to a given value
 */
export type Equals = {
  kind?: '=';
  field: DataField;
  value: AnyType;
};

