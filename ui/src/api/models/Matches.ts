/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { DataField } from './DataField';

/**
 * Match if the field of an alert matches a given regex. The regex can optionally
 * start and end with `/`, to make it easier to spot in configuration. The regex method
 * used is a search (`re.search`), so for strict matches, the `^` and `$` need to be used.
 */
export type Matches = {
  kind?: 'matches';
  field: DataField;
  /**
   * The pattern of the regex. Need to be a valid regex.
   */
  value: string;
  /**
   * Activate the option to ignore the case
   */
  ignore_case?: boolean;
};

