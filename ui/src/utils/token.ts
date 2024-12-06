import axios from 'axios'
import { AxiosError } from 'axios'
import type { InternalAxiosRequestConfig } from 'axios'

import { router } from '@/router'

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
