import axios from "axios"
import type { ListOf, Pagination, TimeRangeParams } from "@/api"
import type { Tag, Group } from '@/api'

export type Snooze = {
  id?: string

  groups: Group[]

  reason: string
  tags: Tag[]

  cancelled?: SnoozeCancel

  startsAt: number
  endsAt: number

  username?: string
}

export type SnoozeCancel = {
  by: string
  at: number
  reason: string
}

export type GetSnoozesParams = {
  timerange: TimeRangeParams
  pagination: Pagination
  search?: string
  filter: string
}

export function getSnoozes(params: GetSnoozesParams): Promise<ListOf<Snooze>> {
  var q = {
    start: params.timerange.start,
    end: params.timerange.end,
    page: params.pagination.page,
    size: params.pagination.pageSize,
    search: params.search,
    filter: params.filter,
  }
  return axios.get<ListOf<Snooze>>("/api/snoozes", {params: q})
    .then((resp) => {
      return resp.data
    })
}

export function createSnooze(item: Snooze): Promise<void> {
  console.log(`createSnooze(${JSON.stringify(item)})`)
  return axios.post("/api/snooze", item)
}

export type PostSnoozeCancelParams = {
  ids: string[]
  reason: string
}

export function cancelSnooze(ids: string[], reason: string): Promise<void> {
  console.log(`cancelSnooze(${JSON.stringify(ids)}, '${reason}')`)
  var query = {
    ids: ids,
    reason: reason,
  }
  return axios.post("/api/snooze/cancel", query)
}
