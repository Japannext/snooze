import { default as axios, AxiosInstance, AxiosRequestConfig } from 'axios'
import { ConditionObject } from '@/utils/condition'

import {DatabaseItem} from '@/utils/types'

function tokenInterceptor(config: AxiosRequestConfig) {
  const token = localStorage.getItem('snooze-token')
  if (token) {
    config.headers['Authorization'] = `JWT ${token}`
  }
  return config
}

export class Api {
  axios: AxiosInstance
  constructor(baseURL: string) {
    this.axios = axios.create({baseURL: baseURL})
    this.axios.defaults.headers.post['Content-type'] = 'application/json'
    this.axios.interceptors.request.use(tokenInterceptor, (err) => Promise.reject(err))
  }
  endpoint(collection: string): Endpoint {
    const endpoint = new Endpoint(this, collection)
    return endpoint
  }
}

class Endpoint {
  axios: AxiosInstance
  url: string
  constructor(api: Api, collection: string) {
    this.axios = api.axios
    this.url = `/${collection}`
  }

  insert_one(item: object): Promise<databaseItem> {
    return this.insert_many([item])
    .then(changeResult => {
      if (changeResult.added && changeResult.added.length > 0) {
        return changeResult.added[0]
      } else if (changeResult.rejected && changeResult.rejected.length > 0) {
        const rejected: RejectedItem = changeResult.rejected[0]
        throw rejected.error
      } else {
        throw `No data added or rejected: ${changeResult}`
      }
    })
  }

  insert_many(items: object[]): Promise<ChangeResult> {
    const itemsToAdd = items.map(item => filterObject(item))
    return this.axios.post(this.url, itemsToAdd)
    .then(response => {
      const result: Result<ChangeResult> = response.data
      return result.data
    })
    .catch(errorHandler)
  }

  find(query: Query, options: PaginationOptions): Promise<databaseItem[]> {
    const params = {
      s: query,
      ...options,
    }
    //log.info(`GET ${this.url}`, params)
    return this.axios.get(this.url, {params: params})
    .then(response => {
      const result: Result<databaseItem[]> = response.data
      const data = result.data
      //log.info('results', data)
      return data
    })
    .catch(errorHandler)
  }

  update_one(item: databaseItem): Promise<databaseItem> {
    return this.update_many([item])
    .then(changeResult => {
      if (changeResult.updated && changeResult.updated.length > 0) {
        return changeResult.updated[0]
      } else if (changeResult.rejected && changeResult.rejected.length > 0) {
        const rejected: RejectedItem = changeResult.rejected[0]
        throw rejected['error']
      } else {
        throw `No data updated or rejected: ${changeResult}`
      }
    })
    .catch(errorHandler)
  }

  update_many(items: databaseItem[]): Promise<ChangeResult> {
    const itemsToUpdate = items.map(item => filterObject(item))
    return this.axios.put(this.url, itemsToUpdate)
    .then(response => {
      const result: Result<ChangeResult> = response.data
      return result.data
    })
    .catch(errorHandler)
  }

  delete_many(uids: string[]) {
    const uidConditions = uids.map(uid => new ConditionObject('=', ['uid', uid]))
    const query = new ConditionObject('OR', uidConditions)
    const params = {
      s: query.toArray(),
    }
    this.axios.delete(this.url, {params: params})
    .then(response => {
      const result: Result<databaseItem[]> = response.data
      if (result.count != uids.length) {
        throw `Deleted an unexpected number of items: ${result.count} instead of ${uids.length}`
      }
    })
    .catch(errorHandler)
  }
}



function filterObject(item: object): object {
  const asArray = Object.entries(item)
  const filteredArray = asArray.filter(([key, ]) => key[0] == '_')
  return Object.fromEntries(filteredArray)
}

function errorHandler(error: Error): never {
  if (axios.isAxiosError(error)) {
    // Type error
  }
  //log.error(error)
  throw error
}
