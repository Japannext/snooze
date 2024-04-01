/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type ValidationError = {
  name?: string;
  text: string;
  field_errors: Record<string, Array<string>>;
};

