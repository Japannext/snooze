import axios from "axios"
import type { ListOf, Pagination, TimeRangeParams } from "@/api"

export type Alert = {
  id?: string;

  status: object;
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
