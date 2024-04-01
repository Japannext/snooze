import type { Direction, CancelablePromise, Body_update } from '@/api'

export interface SearchParams {
    query?: string,
    pageNb?: number,
    perPage?: number,
    orderBy?: string,
    direction?: Direction,
}

export interface Collection<T> {
  items: Array<T>,
}

export interface CrudService<T> {
  search(
    pageNb?: number,
    perPage?: number,
    orderBy?: string,
    direction?: Direction,
  ): CancelablePromise<Collection<T>>

  get(oid: string): CancelablePromise<T>
  create(requestBody: T): CancelablePromise<any>
  update(requestBody: Body_update): CancelablePromise<any>
  delete(requestBody: Array<string>): CancelablePromise<any>
  validate(requestBody: T): CancelablePromise<any>

}
