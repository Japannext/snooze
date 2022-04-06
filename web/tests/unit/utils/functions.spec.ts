import { trimDate } from '@/utils/functions'
import moment from 'moment'
import MockDate from 'mockdate'

// Necessary because moment throw warnings during the 'invalid time' test
moment.suppressDeprecationWarnings = true

describe('trimDate', () => {
  test('same day', () => {
    MockDate.set('2022-01-10T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(trimDate(date, false)).toBe('Today 10:00')
  })
  test('same day (with seconds)', () => {
    MockDate.set('2022-01-10T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(trimDate(date, true)).toBe('Today 10:00:00')
  })
  test('different day', () => {
    MockDate.set('2022-01-12T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(trimDate(date, false)).toBe('Jan 10th 10:00')
  })
  test('different day (with seconds)', () => {
    MockDate.set('2022-01-12T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(trimDate(date, true)).toBe('Jan 10th 10:00:00')
  })
  test('different year', () => {
    MockDate.set('2023-01-12T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(trimDate(date, false)).toBe('Jan 10th 2022')
  })
  test('different year (with seconds)', () => {
    MockDate.set('2023-01-12T12:02:30+09:00')
    const date = '2022-01-10T10:00:00+09:00'
    expect(trimDate(date, true)).toBe('Jan 10th 2022')
  })
  test('empty time', () => {
    const date = ''
    expect(trimDate(date, false)).toBe('Empty')
  })
  test('invalid time', () => {
    const date = 'asdfasdlkfj'
    expect(trimDate(date, false)).toBe('Invalid date')
  })
})
