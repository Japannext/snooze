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

export function createTag(item: Tag): Promise<void> {
  return axios.post("/api/tag", item)
}
