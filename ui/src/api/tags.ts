import axios from "axios"
import type { ListOf, Pagination, TimeRangeParams } from "@/api"

export type Tag = {
  id?: string

  name: string
  description: string
  color: string
}

export type GetTagsParams = {
  timerange: TimeRangeParams,
  pagination: Pagination
  search?: string
}

export function getTags(params: GetTagsParams): Promise<ListOf<Tag>> {
  var q = {
    search: params.search,
  }
  return axios.get<ListOf<Tag>>("/api/tags", {params: q})
    .then((resp) => {
      return resp.data
    })
}

export function createTag(item: Tag): Promise {
  return axios.post("/api/tag", item)
}
