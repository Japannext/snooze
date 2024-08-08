import axios from 'axios'
import type { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';

import type { Pagination } from './pagination';

export type Alert = {
  source: Source;

  timestamp?: number;
  observedTimestamp?: number;

  groupHash?: string;
  groupLabels?: any;

  severityText: string;
  severityNumber: number;

  labels?: any;
  attributes?: any;
  body?: any;

  mute?: Mute;
}

export type Source = {
  kind: string
  name?: string
}

export type Mute = {
  enabled: boolean;
  component: string;
  rule: string;
  skipNotification: boolean;
  skipStorage: boolean;
  silentTest: boolean;
  humanTest: boolean;
}

export type AlertList = {
  items: Array<Alert>
}

export function searchAlerts(query: string, pagination: Pagination): Promise<Array<Alert>> {
  return axios.get<AlertList>('/v2/alerts', {params: {query: query, ...pagination}})
    .then(function (resp: AxiosResponse<AlertList, any>): Array<Alert> {
      return resp.items
    })
    .catch(function (error) {
      console.log(error)
    })
}
