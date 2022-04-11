// Functions related to condition transformation

import { uuid } from 'vue-uuid'

export const OPERATION_TYPE = {
  '=': 'binary',
  '!=': 'binary',
  '>': 'binary',
  '>=': 'binary',
  '<': 'binary',
  '<=': 'binary',
  'MATCHES': 'binary',
  'CONTAINS': 'binary',
  'EXISTS': 'unary',
  'SEARCH': 'unary',
  'AND': 'logic',
  'OR': 'logic',
  'NOT': 'not',
  '': 'alwaysTrue',
  null: 'alwaysTrue',
  undefined: 'alwaysTrue',
}

export const OPERATION_SYMBOL = {
  'MATCHES': '~',
  'SEARCH': 'unary',
  'AND': '&',
  'OR': '|',
  'NOT': '!',
}

// Class representing a condition
export class ConditionObject {
  /** Create a condition
   * @param {string} op The condition operator
   * @param {string[]|ConditionObject[]} args The arguments passed to the condition operator
  **/
  constructor(op, args) {
    this.id = uuid.v4()
    this.args = args
    this.operation = op
  }
  /** Create a condition from a nested array
   * @param array A nested array representing a condition
   * @returns {ConditionObject}
  **/
  static fromArray(array) {
    if (array === undefined) {
      return new ConditionObject('', [])
    }
    var operation = array[0]
    let args
    if (OPERATION_TYPE[operation] == 'logic' || operation == 'NOT') {
      args = array.slice(1).map(arg => ConditionObject.fromArray(arg))
    } else {
      args =  array.slice(1)
    }
    var c = new ConditionObject(operation, args)
    return c
  }
  get type() {
    return OPERATION_TYPE[this.operation]
  }
  get operationSymbol() {
    return this.operation in OPERATION_SYMBOL ? OPERATION_SYMBOL[this.operation] : this.operation
  }
  toJSON() {
    var json = {operation: this.operation, id: this.id}
    if (this.type == 'logic') {
      json['args'] = this.args.map(arg => arg.toJSON())
    } else {
      json['args'] = this.args
    }
    return json
  }
  /** An array representation of the condition
   * @returns
  **/
  toArray() {
    if (this.type == 'logic' || this.type == 'not') {
      return [this.operation].concat(this.args.map(arg => arg.toArray()))
    } else {
      return [this.operation].concat(this.args)
    }
  }
  /** Return a query language representation of the condition.
   * @returns {string} A query language representation of the condition
  **/
  toSearch() {
    switch(this.type) {
      case 'logic':
        return "(" + this.args.map(arg => arg.toSearch()).join(` ${this.operationSymbol} `) + ")"
      case 'binary':
        return `${this.args[0]} ${this.operationSymbol} ${JSON.stringify(this.args[1])}`
      case 'unary':
        return `(${this.operation} ${JSON.stringify(this.args[0])})`
      case 'not':
        return `(!${this.args[0]})`
      case 'alwaysTrue':
        return '()'
      default:
        return `Invalid condition`
    }
  }
  /** Returns an HTML representation of the condition.
   * @returns {string}
  **/
  toHTML() {
    switch(this.type) {
      case 'logic':
        return "(" + this.args.map(arg => arg.toHTML()).join(` <b>${this.operation}</b> `) + ")"
      case 'binary':
        return `${this.args[0]} <b>${this.operationSymbol}</b> ${JSON.stringify(this.args[1])}`
      case 'unary':
        return `(<b>${this.operation}</b> ${JSON.stringify(this.args[0])})`
      case 'not':
        return `<b>NOT</b> ${this.args[0].toHTML()}`
      case 'alwaysTrue':
        return '<b>Always true</b>'
      default:
        return `<b>Invalid Condition</b>`
    }
  }
  /** Combine two valid condition with an operation (AND, OR)
   * @param {(AND|OR)} operator The logic operator to use to combine
   * @param {ConditionObject} other The other condition to combine with
   * @returns {ConditionObject} The result of the combination
  **/
  combine(operator, others) {
    return new ConditionObject(operator, [this].concat(others))
  }
  /** A string representation of the condition (a JSON stringify view). Used for debug logs.
   * @returns {string}
  **/
  toString() {
    return JSON.stringify(this.toJSON())
  }

}
