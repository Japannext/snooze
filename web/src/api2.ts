import jwt_decode from 'jwt-decode'
import { is } from 'typescript-is'

import { Api } from '@/utils/api2'

let url: string
if ('VUE_APP_API' in process.env) {
  url = process.env.VUE_APP_API
} else {
  url = '/api/'
}

export const api2 = new Api(url)

type UserPayload = {
  user: string
  method: string
}

if (!localStorage.getItem('username')) {
  const token = localStorage.getItem('snooze-token')
  if (token !== null) {
    const payload = jwt_decode(token)
    if (is<UserPayload>(payload)) {
      localStorage.setItem('username', payload.user)
    }
  }
}
