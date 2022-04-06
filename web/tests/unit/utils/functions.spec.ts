import { prettyDate, prettyDuration, prettyNumber, capitalizeFirstLetter } from '@/utils/functions'
import moment from 'moment'
import MockDate from 'mockdate'

// Necessary because moment throw warnings during the 'invalid time' test
moment.suppressDeprecationWarnings = true

describe('prettyDate', () => {
  test('same day', () => {
    MockDate.set('2022-01-10T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(prettyDate(date, false)).toBe('Today 10:00')
  })
  test('same day (with seconds)', () => {
    MockDate.set('2022-01-10T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(prettyDate(date, true)).toBe('Today 10:00:00')
  })
  test('different day', () => {
    MockDate.set('2022-01-12T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(prettyDate(date, false)).toBe('Jan 10th 10:00')
  })
  test('different day (with seconds)', () => {
    MockDate.set('2022-01-12T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(prettyDate(date, true)).toBe('Jan 10th 10:00:00')
  })
  test('different year', () => {
    MockDate.set('2023-01-12T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(prettyDate(date, false)).toBe('Jan 10th 2022')
  })
  test('different year (with seconds)', () => {
    MockDate.set('2023-01-12T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(prettyDate(date, true)).toBe('Jan 10th 2022')
  })
  test('empty time', () => {
    const date = ''
    expect(prettyDate(date, false)).toBe('Empty')
  })
  test('invalid time', () => {
    const date = 'asdfasdlkfj'
    expect(prettyDate(date, false)).toBe('Invalid date')
  })
})

describe('prettyDuration', () => {
  test('Number input', () => {
    expect(prettyDuration(4000)).toBe('1h 6m 40s')
  })
  test('String input', () => {
    expect(prettyDuration('5000')).toBe('1h 23m 20s')
  })
})

describe('prettyNumber', () => {
  test('Small number', () => {
    expect(prettyNumber(42)).toBe('42')
  })
  test('Big number', () => {
    expect(prettyNumber(4200000)).toBe('4,200,000')
  })
})

describe('capitalizeFirstLetter', () => {
  test('test', () => {
    expect(capitalizeFirstLetter('title')).toBe('Title')
  })
})
