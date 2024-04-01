/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { DataField } from './DataField';

/**
 * Match if a given field exist and is not null in the alert
 */
export type Exists = {
  kind?: 'exists';
  field: DataField;
};

