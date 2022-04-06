import MockAdapter from 'axios-mock-adapter'

import { Api } from '@/utils/api2'

const api = new Api('http://example.com')
const mock = new MockAdapter(api.axios)

describe('Endpoint', () => {
  const snooze = api.endpoint('snooze')
  describe('find', () => {
    it('full match', () => {
      const snoozeFilters = [
        {},
        {},
      ]
      mock.onGet('/api/snooze', {s: []})
        .reply(200, snoozeFilters)
      snooze.find([])
        .then(results => {
          expect(results.length).toBe(2)
        })
    })
  })
})
