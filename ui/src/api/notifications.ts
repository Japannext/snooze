import axios from 'axios'
import type { Destination } from '@/api'
import type { ListOf, Pagination, TimeRangeParams, } from '@/api'

export type Notification = {
  _id: string
  notificationTime: number
  destination: Destination
  acknowledged: boolean

  type: string
  itemID: string
  source: object
  identity: object
  message: string
  labels: object
  documentationURL: string
}

export type GetNotificationsParams = {
  timerange: TimeRangeParams
  pagination: Pagination
  search?: string
  filter?: string
}

export function getNotifications(params: GetNotificationsParams): Promise<ListOf<Notification>> {
  var q = {
    start: params.timerange.start,
    end: params.timerange.end,
    page: params.pagination.page,
    size: params.pagination.pageSize,
    search: params.search,
    filter: params.filter,
  }
  return axios.get<ListOf<Notification>>("/api/notifications", {params: q})
    .then((resp) => {
      return resp.data
    })
}
