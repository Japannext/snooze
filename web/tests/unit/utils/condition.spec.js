import { ConditionObject } from '@/utils/condition'

describe('ConditionObject', () => {

  describe("fromArray", () => {
    test("a=x", () => {
      var condition = ConditionObject.fromArray(['=', 'a', 'x'])
      expect(condition.operation).toBe("=")
      expect(condition.args).toStrictEqual(['a', 'x'])
    })
    test("a=x & b=y", () => {
      var condition = ConditionObject.fromArray(['AND', ['=', 'a', 'x'], ['=', 'b', 'y']])
      expect(condition.operation).toBe('AND')
      expect(condition.args[0].operation).toBe('=')
      expect(condition.args[0].args).toStrictEqual(['a', 'x'])
      expect(condition.args[1].operation).toBe('=')
      expect(condition.args[1].args).toStrictEqual(['b', 'y'])
    })
  })

  describe("toArray", () => {
    test("a=x", () => {
      var condition = new ConditionObject('=', ['a', 'x'])
      expect(condition.toArray()).toStrictEqual(['=', 'a', 'x'])
    })
    test("a=x & b=y", () => {
      var cond1 = new ConditionObject('=', ['a', 'x'])
      var cond2 = new ConditionObject('=', ['b', 'y'])
      var condition = new ConditionObject('AND', [cond1, cond2])
      expect(condition.toArray()).toStrictEqual(['AND', ['=', 'a', 'x'], ['=', 'b', 'y']])
    })
  })

  describe("toSearch", () => {
    test("a=x", () => {
      var condition = new ConditionObject('=', ['a', 'x'])
      expect(condition.toSearch()).toBe('a = "x"')
    })
    test("a=x & b=y", () => {
      var cond1 = new ConditionObject('=', ['a', 'x'])
      var cond2 = new ConditionObject('=', ['b', 'y'])
      var condition = new ConditionObject('AND', [cond1, cond2])
      expect(condition.toSearch()).toBe('(a = "x" & b = "y")')
    })
  })

  describe("combine", () => {
    test("a=x & b=y", () => {
      var cond1 = new ConditionObject('=', ['a', 'x'])
      var cond2 = new ConditionObject('=', ['b', 'y'])
      var condition = cond1.combine('AND', cond2)
      expect(condition.operation).toBe('AND')
      expect(condition.args[0].operation).toBe('=')
      expect(condition.args[0].args).toStrictEqual(['a', 'x'])
      expect(condition.args[1].operation).toBe('=')
      expect(condition.args[1].args).toStrictEqual(['b', 'y'])
    })
  })

})
