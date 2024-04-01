/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { BooleanType } from './BooleanType';
import type { FloatType } from './FloatType';
import type { IntegerType } from './IntegerType';
import type { NullType } from './NullType';
import type { StringType } from './StringType';

export type AnyType = (StringType | IntegerType | FloatType | BooleanType | NullType);

