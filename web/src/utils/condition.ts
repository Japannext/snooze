// Functions related to condition transformation

import { uuid } from 'vue-uuid'
import { is, assertType } from 'typescript-is'

type Unary = 'EXISTS'|'SEARCH'
type Binary = '='|'!='|'>'|'<'|'>='|'<='|'MATCHES'|'CONTAINS'|'IN'
type Logic = 'AND'|'OR'
type AllLogic = Logic|'NOT'
type AlwaysTrue = ''|null|undefined
type AllOperations = Unary|Binary|AlwaysTrue
type AllCondition = AllOperations|AllLogic

type ConditionItem = string | ConditionArray
export type ConditionArray = [AllCondition, ...ConditionItem[]]

const SYMBOLS = new Map<string, string>([
  ['MATCHES', '~'],
  ['AND', '&'],
  ['OR', '|'],
  ['NOT', '!'],
])

// Workaround because typescript-is doesn't support nested types
// properly. Equivalent to is<ConditionArray>(array)
function isConditionArray(array: any): array is ConditionArray {
  return (
    Array.isArray(array) &&
    array.length > 0 &&
    is<AllCondition>(array[0]) &&
    (
      is<string[]>(array) ||
      array.slice(1).every(arg => isConditionArray(arg))
    )
  )
}

// Equivalent to assertType<ConditionArray[]>(array)
function assertConditionArrays(array: any): ConditionArray[] {
  if (array.every((arg: any) => isConditionArray(arg))) {
    return array as ConditionArray[]
  } else {
    throw `Invalid condition: ${array}`
  }
}

export abstract class ConditionAbstract {
  uuid: string
  abstract readonly args: string[] | Array<SingleCondition|MultiCondition>
  abstract readonly operation: AllCondition

  constructor() {
    this.uuid = uuid.v4()
  }

  /** An array representation of the condition
   * @returns {ConditionArray} An array representation of the condition
  **/
  abstract toArray(): ConditionArray

  /** Return a query language representation of the condition.
   * @returns {string} A query language representation of the condition
  **/
  abstract toSearch(): string

  /** Returns an HTML representation of the condition.
   * @returns {string} An HTML representation of a condition
  **/
  abstract toHTML(): string

  abstract toJSON(): object

  static fromArray(array: ConditionArray): SingleCondition|MultiCondition {
    if (array.length <= 1) {
      return new SingleCondition('', []) // AlwaysTrue
    }
    const [operation, ...args2] = array
    if (is<AllOperations>(operation)) {
      return new SingleCondition(operation, assertType<string[]>(args2))
    }
    if (is<AllLogic>(operation)) {
      const args3 = assertConditionArrays(args2).map(arg => ConditionAbstract.fromArray(arg))
      return new MultiCondition(operation, args3)
    }
    const exhaustiveCheck: never = operation
    throw `Unsupported operator: ${exhaustiveCheck}`
  }

  toString(): string {
    return JSON.stringify(this.toJSON())
  }

}

/** Combine two valid condition with an operation (AND, OR)
 * @param {(AND|OR)} operator The logic operator to use to combine
**/
export function combine(operator: 'AND'|'OR', ...conditions: Array<SingleCondition|MultiCondition>): MultiCondition {
  return new MultiCondition(operator, conditions)
}

export class SingleCondition extends ConditionAbstract {
  readonly operation: AllOperations
  readonly args: string[]

  constructor(op: AllOperations, args: string[]) {
    super()
    this.operation = op
    this.args = args
  }

  toArray(): ConditionArray {
    const args2: ConditionItem[] = this.args
    return [this.operation, ...args2]
  }

  toJSON(): object {
    return {
      uuid: this.uuid,
      operation: this.operation,
      args: this.args,
    }
  }

  toSearch(): string {
    if (is<Binary>(this.operation)) {
      const symbol: string = SYMBOLS.get(this.operation) || this.operation
      return `${this.args[0]} ${symbol} ${JSON.stringify(this.args[1])}`
    }
    else if (is<Unary>(this.operation)) {
      return `(${this.operation} ${JSON.stringify(this.args[0])})`
    }
    else if (is<AlwaysTrue>(this.operation)) {
      return '()'
    }
    const exhaustiveCheck: never = this.operation
    throw `Unsupported operator: ${exhaustiveCheck}`
  }

  toHTML(): string {
    if (is<Binary>(this.operation)) {
      const symbol = SYMBOLS.get(this.operation) || this.operation
      return `${this.args[0]} <b>${symbol}</b> ${JSON.stringify(this.args[1])}`
    }
    else if (is<Unary>(this.operation)) {
      return `(<b>${this.operation}</b> ${JSON.stringify(this.args[0])})`
    }
    else if (is<AlwaysTrue>(this.operation)) {
      return '<b>Always true</b>'
    }
    const exhaustiveCheck: never = this.operation
    throw `Unsupported operator: ${exhaustiveCheck}`
  }
}

export class MultiCondition extends ConditionAbstract {
  readonly operation: AllLogic
  readonly args: Array<SingleCondition|MultiCondition>

  constructor(op: Logic|'NOT', args: Array<SingleCondition|MultiCondition>) {
    super()
    this.operation = op
    this.args = args
  }

  toArray(): ConditionArray {
    return [this.operation, ...this.args.map(arg => arg.toArray())]
  }

  toJSON(): object {
    return {
      uuid: this.uuid,
      operation: this.operation,
      args: this.args.map(arg => arg.toJSON()),
    }
  }

  toSearch(): string {
    if (is<Logic>(this.operation)) {
      const symbol = SYMBOLS.get(this.operation) || this.operation
      return "(" + this.args.map(arg => arg.toSearch()).join(` ${symbol} `) + ")"
    }

    else if (this.operation == 'NOT') {
      return `(!${this.args[0]})`
    }

    const exhaustiveCheck: never = this.operation
    throw `Unsupported type ${exhaustiveCheck}`
  }

  toHTML(): string {
    if (is<Logic>(this.operation)) {
      return "(" + this.args.map(arg => arg.toHTML()).join(` <b>${this.operation}</b> `) + ")"
    }

    else if (this.operation == 'NOT') {
      return `<b>NOT</b> ${this.args[0].toHTML()}`
    }

    const exhaustiveCheck: never = this.operation
    throw `Unsupported type ${exhaustiveCheck}`
  }

}
