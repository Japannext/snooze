import axios from 'axios'
import { AxiosError } from 'axios'
import type { InternalAxiosRequestConfig } from 'axios'

import { router } from '@/router'
import { xSnoozeToken } from '@/stores'

// Inject the token
axios.interceptors.request.use(
  (request: InternalAxiosRequestConfig<any>) => {
    if (xSnoozeToken.value != "") {
      request.headers.set('X-Snooze-Token', xSnoozeToken.value)
    } else {
      console.log("Could not find X-Snooze-Token in header: redirecting to /login route")
      router.push({name: 'login'})
    }
    return request
  },
  (error) => { return Promise.reject(error) }

)

// Redirect 401 to /login
axios.interceptors.response.use(
  (resp) => resp,
  (err: Error | AxiosError) => {
    console.log(JSON.stringify(err, null, 2))
    if (
      (err instanceof AxiosError)
      && (err.response !== undefined)
      && (err.response.status == 401)
      && (router.currentRoute.value.path) != '/login'
    ) {
      console.log("Got 401 status code: redirecting to /login route")
      router.push({name: "login"})
    } else {
      return Promise.reject(err)
    }
  },
)
