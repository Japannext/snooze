import axios from "axios"
import type { AxiosResponse } from "axios"
import { DateTime } from "luxon"
import type { ListOf } from "@/api/utils"

export type Log<DT = DateTime> = {
  id?: string;

  source: object;
  identity?: object
  labels?: object;

  actualTime: DT;
  observedTime: DT;
  displayTime: DT;

  groupHash?: string;
  groupLabels?: any;

  severityText: string;
  severityNumber: number;

  message: string;

  mute?: object;
}

type rawLog = Log<number>

export type LogParams = {
  start: number;
  end: number;
  pagination: object;
  search: string;
}

function convertTimes(items: Log<number>[]): Log[] {
  return items.map((item) => {
    return { ...item,
      actualTime: DateTime.fromMillis(item.actualTime),
      displayTime: DateTime.fromMillis(item.displayTime),
      observedTime: DateTime.fromMillis(item.observedTime),
    }
  })
}

export function getLogs(start: number, end: number, pagination: object, search: number): Promise<ListOf<Log>> {
  var params = {
    page: pagination.page,
    size: pagination.pageSize,
    start: start,
    end: end,
    search: search,
  }
  return axios.get<ListOf<rawLog>>("/api/logs", {params: params})
    .then()
    .then(resp => {
      return {...resp.data,
        items: convertTimes(resp.data.items)
      }
    })
}
