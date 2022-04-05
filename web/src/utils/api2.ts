import { default as axios, AxiosInstance } from 'axios'
import { ConditionObject } from '@/utils/condition'

type QueryElement = string | Query[]
type Query = Array<QueryElement>

export class Api {
  axios: AxiosInstance
  constructor(baseURL: string) {
    this.axios = axios.create({baseURL: baseURL})
  }
}

interface RejectedItem {
  error: string
}

interface ChangeResult {
  added?: Object[]
  updated?: Object[]
  replaced?: Object[]
  rejected?: RejectedItem[]
}

interface Result<U> {
  data: U
  count: number
}

export class Endpoint {
  axios: AxiosInstance
  url: string
  constructor(api: Api, collection: string) {
    this.axios = api.axios
    this.url = `${collection}`
  }

  insert_one(item: Object): Promise<Object> {
    return this.insert_many([item])
    .then(changeResult => {
      if (changeResult.added && changeResult.added.length > 0) {
        return changeResult.added[0]
      } else if (changeResult.rejected && changeResult.rejected.length > 0) {
        let rejected: RejectedItem = changeResult.rejected[0]
        throw rejected['error']
      } else {
        throw `No data added or rejected: ${changeResult}`
      }
    })
  }

  insert_many(items: Object[]): Promise<ChangeResult> {
    let itemsToAdd = items.map(item => filterObject(item))
    return this.axios.post(this.url, itemsToAdd)
    .then(response => {
      let result: Result<ChangeResult> = response.data
      return result.data
    })
    .catch(errorHandler)
  }

  find(query: Query): Promise<Object[]> {
    var params = {s: query}
    //log.info(`GET ${this.url}`, params)
    return this.axios.get(this.url, {params: params})
    .then(response => {
      const result: Result<Object[]> = response.data
      const data = result.data
      //log.info('results', data)
      return data
    })
    .catch(errorHandler)
  }

  update_one(uid: string, item: Object): Promise<Object> {
    let itemToUpdate: Object = {...item, uid: uid}
    return this.update_many([itemToUpdate])
    .then(changeResult => {
      if (changeResult.updated && changeResult.updated.length > 0) {
        return changeResult.updated[0]
      } else if (changeResult.rejected && changeResult.rejected.length > 0) {
        let rejected: RejectedItem = changeResult.rejected[0]
        throw rejected['error']
      } else {
        throw `No data updated or rejected: ${changeResult}`
      }
    })
    .catch(errorHandler)
  }

  update_many(items: Object[]): Promise<ChangeResult> {
    let itemsToUpdate = items.map(item => filterObject(item))
    return this.axios.put(this.url, itemsToUpdate)
    .then(response => {
      let result: Result<ChangeResult> = response.data
      return result.data
    })
    .catch(errorHandler)
  }

  delete_many(uids: string[]) {
    let uidConditions = uids.map(uid => new ConditionObject('=', ['uid', uid]))
    let query = new ConditionObject('OR', uidConditions)
    let params = {
      s: query.toArray(),
    }
    this.axios.delete(this.url, {params: params})
    .then(response => {
      let result: Result<Object[]> = response.data
      if (result.count != uids.length) {
        throw `Deleted an unexpected number of items: ${result.count} instead of ${uids.length}`
      }
    })
    .catch(errorHandler)
  }
}

function filterObject(item: Object): Object {
  const asArray = Object.entries(item)
  const filteredArray = asArray.filter(([key, value]) => key[0] == '_')
  return Object.fromEntries(filteredArray)
}

function errorHandler(error: Error): never {
  if (axios.isAxiosError(error)) {
    // Type error
  } else {
  }
  //log.error(error)
  throw error
}
