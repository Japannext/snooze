/**
 * A wrapper around axios return promises to provide more error handling features.
**/

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios'
import { is } from 'typescript-is'

type AxiosSuccessCallback<T, O = T|AxiosResponse<T>|void> = (data: AxiosResponse<T>) => PromiseLike<O>
type SuccessCallback<T, O=T|void> = (data: T) => PromiseLike<O>
type RejectCallback<E> = (error: E) => PromiseLike<never>

export class TypeCheckingError extends Error {
  constructor(obj: object) {
    super(JSON.stringify(obj))
    Object.setPrototypeOf(this, TypeCheckingError.prototype)
  }
}

const HTTP_ERROR_CODES = new Map<string, number>([
  ['Bad Request', 400],
  ['Unauthorized', 401],
  ['Payment Required', 402],
  ['Forbidden', 403],
  ['Not Found', 404],
  ['Method Not Allowed', 405],
  ['Not Acceptable', 406],
  ['Proxy Authentication Required', 407],
  ['Request Timeout', 408],
])





/** Wrapping Axios (it's not possible to extend Axios since they don't expose it)
**/
export class Api {
  axios: AxiosInstance
  constructor(config: AxiosRequestConfig) {
    this.axios = axios.create(config)
  }
  get<T>(url: string, config?: AxiosRequestConfig): HttpPromise<T> {
    return new HttpPromise(this.axios.get(url, config))
  }
  head<T>(url: string, config?: AxiosRequestConfig): HttpPromise<T> {
    return new HttpPromise(this.axios.head(url, config))
  }
  delete<T>(url: string, config?: AxiosRequestConfig): HttpPromise<T> {
    return new HttpPromise(this.axios.delete(url, config))
  }
  options<T>(url: string, config?: AxiosRequestConfig): HttpPromise<T> {
    return new HttpPromise(this.axios.options(url, config))
  }
  post<T, D=any>(url: string, data?: D, config?: AxiosRequestConfig): HttpPromise<T> {
    return new HttpPromise(this.axios.post(url, data, config))
  }
  put<T, D=any>(url: string, data?: D, config?: AxiosRequestConfig): HttpPromise<T> {
    return new HttpPromise(this.axios.put(url, data, config))
  }
  patch<T, D=any>(url: string, data?: D, config?: AxiosRequestConfig): HttpPromise<T> {
    return new HttpPromise(this.axios.patch(url, data, config))
  }
}

/** A promise with more features related to HTTP calls
 * 1) .catch() can take a HTTP code argument to catch only a certain
 *    type of HTTP errors (akin to the behavior of bluebird)
 * 2) Ensure the type checking of .then() at runtime thanks to `typescript-is`'s assertType,
 *    thus ensuring a consistent check
**/
export class HttpPromise<T=any, O=any> {
  codeErrors: Map<number, (err: Error) => never>
  typeError: ((err: Error) => never) | null
  successCallback: (data: T) => any

  constructor(promise: Promise<any>) {
    this.successCallback = (data: T) => {
      return data
    }
    this.codeErrors = new Map()
    this.typeError = null
    promise.then(this.computeSuccess, this.computeErrors)
  }

  private computeSuccess(response: AxiosResponse<T>): O {
    if (is<T>(response.data)) {
      return this.successCallback(response.data)
    } else {
      throw new TypeCheckingError(response)
    }
  }

  private computeErrors(error: Error): never {
    if (this.typeError != null && error instanceof TypeCheckingError) {
      this.typeError(error)
    }
    if (!axios.isAxiosError(error)) {
      throw error
    }
    if (error.response && error.response.status) {
      const callback = this.codeErrors.get(error.response.status)
      if (callback !== undefined) {
        callback(error)
      }
      throw error
    }
    throw error
  }

  /** Execute a callback on success. Guarantee the type with type checking
   * Use parameter overloading to make sure it returns a Promise<void> when passing
   * a callback that doesn't return (useful inside vue or tests)
  **/
  onSuccess(callback: (data: T) => any): this {
    this.successCallback = callback
    return this
  }

  onError(code: number, callback: (err: Error) => never): this {
    this.codeErrors.set(code, callback)
    return this
  }

  /** Execute a callback on typechecking error
  **/
  onTypeError(onRejected: (error: Error) => never): this {
    this.typeError = (error => {
      if (error instanceof TypeCheckingError) {
        onRejected(error)
      } else {
        throw error
      }
    })
    return this
  }

}
