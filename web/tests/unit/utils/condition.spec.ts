import { assertType } from 'typescript-is'
import { ConditionAbstract, SingleCondition, MultiCondition, combine } from '@/utils/condition'

describe('ConditionAbstract', () => {

  describe('fromArray', () => {
    test('a=x', () => {
      const condition = ConditionAbstract.fromArray(['=', 'a', 'x'])
      expect(condition.operation).toBe('=')
      expect(condition.args).toStrictEqual(['a', 'x'])
    })
    test("a=x & b=y", () => {
      const condition = ConditionAbstract.fromArray(['AND', ['=', 'a', 'x'], ['=', 'b', 'y']])
      expect(condition.operation).toBe('AND')
      const cond1 = assertType<SingleCondition>(condition.args[0])
      const cond2 = assertType<SingleCondition>(condition.args[1])
      expect(cond1.operation).toBe('=')
      expect(cond1.args).toStrictEqual(['a', 'x'])
      expect(cond2.operation).toBe('=')
      expect(cond2.args).toStrictEqual(['b', 'y'])
    })

  describe('toSearch', () => {
    test("a=x", () => {
      const condition = new SingleCondition('=', ['a', 'x'])
      expect(condition.toSearch()).toBe('a = "x"')
    })
    test("a=x & b=y", () => {
      const cond1 = new SingleCondition('=', ['a', 'x'])
      const cond2 = new SingleCondition('=', ['b', 'y'])
      const condition = new MultiCondition('AND', [cond1, cond2])
      expect(condition.toSearch()).toBe('(a = "x" & b = "y")')
    })
  })

  describe("combine", () => {
    test("a=x & b=y", () => {
      const cond1 = new SingleCondition('=', ['a', 'x'])
      const cond2 = new SingleCondition('=', ['b', 'y'])
      const condition = combine('AND', cond1, cond2)
      expect(condition.operation).toBe('AND')
      expect(condition.args[0].operation).toBe('=')
      expect(condition.args[0].args).toStrictEqual(['a', 'x'])
      expect(condition.args[1].operation).toBe('=')
      expect(condition.args[1].args).toStrictEqual(['b', 'y'])
    })
  })
  })

})
