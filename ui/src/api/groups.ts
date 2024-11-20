import axios from "axios"
import type { ListOf, Pagination, TimeRangeParams } from "@/api"

export type Group = {
  id?: string;

}

export type GetGroupsParams = {
  search?: string
}

export function getGroups(params: GetGroupsParams): Promise<ListOf<Group>> {
  var q = {
    search: params.search,
  }
  return axios.get<ListOf<Group>>("/api/groups", {params: q})
    .then((resp) => {
      return resp.data
    })
}
