import axios from 'axios'
import { AxiosError } from 'axios'
import type { InternalAxiosRequestConfig } from 'axios'

import { router } from '@/router'
import { ApiError } from '@/api'
import { xSnoozeToken } from '@/stores'

// Inject the token
axios.interceptors.request.use(
  (request: InternalAxiosRequestConfig<any>) => {
    if (xSnoozeToken.value != null) {
      request.headers.set('X-Snooze-Token', xSnoozeToken.value)
    } else {
      router.push({name: 'login'})
    }
    return request
  },
  (error) => { return Promise.reject(error) }

)

// Redirect 401 to /login
axios.interceptors.response.use(
  (resp) => resp,
  (err: Error | ApiError | AxiosError) => {
    console.log(JSON.stringify(err, null, 2))
    if (
      (err instanceof AxiosError)
      && (err.response !== undefined)
      && (err.response.status == 401)
      && (router.currentRoute.value.path) != '/login'
    ) {
      router.push({name: "login"})
    } else {
      return Promise.reject(err)
    }
  },
)
