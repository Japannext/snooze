import MockAdapter from 'axios-mock-adapter'

import { Api, Endpoint } from '@/utils/api2'

const api = new Api('http://example.com')
const mock = new MockAdapter(api.axios)

describe('Endpoint', () => {
  describe('find', () => {
    it('full match', () => {
      let snoozeFilters = [
        {},
        {},
      ]
      let snooze = new Endpoint(api, 'snooze')
      snooze.find([])
        .then(results => {
          expect(results.length).toBe(2)
        })
    })
  })
})
