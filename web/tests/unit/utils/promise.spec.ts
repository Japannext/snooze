/**
 * Tests for the promise library
**/

import MockAdapter from 'axios-mock-adapter'

import { AxiosResponse } from 'axios'
import { Api, HttpPromise, TypeCheckingError } from '@/utils/promise'

const api = new Api({baseURL: 'http://example.com'})
const mock = new MockAdapter(api.axios)

describe('HttpPromise', () => {

  interface MyObject {
    uid: string
    name: string
  }

  interface MyCustomError {
    message: string
  }

  test('onSuccess', () => {
    const objects = [
      {uid: '0', name: 'test 1'},
      {uid: '1', name: 'test 2'},
    ]
    mock.onGet('/api/objects')
      .reply(200, objects)
    api.get<MyObject[]>('/api/objects')
      .onSuccess(data => {
        expect(data.length).toBe(2)
        expect(data[0].uid).toBe('0')
        expect(data[1].uid).toBe('1')
      })
      .onError(404, (error: Error) => {
        throw `404 not found: ${error}`
      })
  })

  test('onSuccess (TypeCheckingError)', () => {
    const objects = [
      {uid: '0', name: 123},
      {uid: '1', name: 456},
    ]
    mock.onGet('/api/objects')
      .reply(200, objects)
    expect(() => {
      api.get<MyObject[]>('/api/objects')
        .onSuccess(data => {})
    }).toThrow(TypeCheckingError)
  })

  test('onErrorCode', () => {
    mock.onGet('/api/secrets')
      .reply(403, {message: "Forbidden!"})

    const test1 = () => {
    api.get('/api/secrets')
      .onError(403, error => {
        throw `custom 403 error`
      })
    }
    expect(test1).toThrow('custom 403 error')

  })

})
