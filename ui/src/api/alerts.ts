import axios from "axios"
import type { ListOf, Pagination, TimeRangeParams, Identity } from "@/api"

export type Alert = {
  id?: string

  hash: string
  source: object
  identity: Identity

  startAt: number
  endAt: number

  alertName: string
  alertGroup: string

  severityText?: string
  severityNumber?: number

  description: string
  summary: string

  labels: object
}

export type GetAlertsParams = {
  timerange: TimeRangeParams,
  pagination: Pagination
  search?: string
  filter?: string
}

export function getAlerts(params: GetAlertsParams): Promise<ListOf<Alert>> {
  var q = {
    start: params.timerange.start,
    end: params.timerange.end,
    page: params.pagination.page,
    size: params.pagination.pageSize,
    search: params.search,
    filter: params.filter,
  }
  return axios.get<ListOf<Alert>>("/api/alerts", {params: q})
    .then((resp) => {
      return resp.data
    })
}
