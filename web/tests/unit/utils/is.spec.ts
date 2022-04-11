import { is, assertType } from 'typescript-is'

test('', () => {
  type Unary = 'EXISTS'|'SEARCH'
  type Binary = '='|'!='|'>'|'<'|'>='|'<='|'MATCHES'|'CONTAINS'|'IN'
  type Logic = 'AND'|'OR'
  type AllLogic = Logic|'NOT'
  type AlwaysTrue = ''|null|undefined
  type AllOperations = Unary|Binary|AlwaysTrue
  type AllCondition = AllOperations|AllLogic

  type Child = string | Nested
  type Nested = [string, ...Child[]]
  //is<Nested>(['a', 'b', 'c'])

  class Test {
  }
  const t = new Test()
  function f(t: Test) {
    console.log(`${t}`)
  }

  f(t)

  assertType<string[]>(['a', 'b', 'c'])
  is<AllCondition>('any')
  is<Logic>('test')
  is<Unary>('test')
  is<AllLogic>('test')
  is<AllOperations>('test')
  is<AlwaysTrue>(null)
  is<AlwaysTrue>(undefined)
  is<AlwaysTrue>('')
  is<Binary>('AND')
  is<'NOT'>('NOT')

})
