import axios from "axios"
import type { AxiosResponse } from "axios"
import type { ListOf, Pagination, TimeRangeParams } from "@/api"

export type Log = {
  id?: string;

  source: object;
  identity?: object
  labels?: object;

  actualTime: number;
  observedTime: number;
  displayTime: number;

  groupHash?: string;
  groupLabels?: any;

  severityText: string;
  severityNumber: number;

  message: string;

  status?: object;
}

export type GetLogsParams = {
  timerange: TimeRangeParams,
  pagination: Pagination
  search?: string
  filter?: string
}

export function getLogs(params: GetLogsParams): Promise<ListOf<Log>> {
  var q = {
    start: params.timerange.start,
    end: params.timerange.end,
    page: params.pagination.page,
    size: params.pagination.pageSize,
    search: params.search,
    filter: params.filter,
  }
  return axios.get<ListOf<Log>>("/api/logs", {params: q})
    .then((resp) => {
      return resp.data
    })
}
