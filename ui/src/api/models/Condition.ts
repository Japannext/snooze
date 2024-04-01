/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { AlwaysTrue } from './AlwaysTrue';
import type { And } from './And';
import type { Equals } from './Equals';
import type { Exists } from './Exists';
import type { GreaterOrEqual } from './GreaterOrEqual';
import type { GreaterThan } from './GreaterThan';
import type { LowerOrEqual } from './LowerOrEqual';
import type { LowerThan } from './LowerThan';
import type { Matches } from './Matches';
import type { Not } from './Not';
import type { NotEquals } from './NotEquals';
import type { Or } from './Or';

/**
 * A condition
 */
export type Condition = (AlwaysTrue | Equals | NotEquals | GreaterThan | LowerThan | GreaterOrEqual | LowerOrEqual | Matches | Exists | And | Or | Not);

