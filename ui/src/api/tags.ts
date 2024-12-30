import axios from "axios"
import type { ListOf, Pagination, TimeRangeParams } from "@/api"

export type Tag = {
  id?: string

  name: string
  description: string
  color: string
}

type TagColor = "default" | "success" | "error" | "warning" | "primary" | "info" | undefined

export type GetTagsParams = {
  search?: string
  pagination: Pagination
}

export type TagList = ListOf<Tag>

export function getTags(params: GetTagsParams): Promise<TagList> {
  var q = {
    search: params.search,
    page: params.pagination.page,
    size: params.pagination.pageSize,
  }
  return axios.get<ListOf<Tag>>("/api/tags", {params: q})
    .then((resp) => {
      return resp.data
    })
}

export function createTag(item: Tag): Promise<void> {
  return axios.post("/api/tag", item)
}
