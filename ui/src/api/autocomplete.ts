import axios from "axios";
import type { ListOf, Pagination, TimeRangeParams, Identity } from "@/api";

export type Log = {
  id?: string;

  source: object;
  identity: Identity;
  labels?: object;

  actualTime: number;
  observedTime: number;
  displayTime: number;

  groupHash?: string;
  groupLabels?: any;

  pattern?: string;
  profile?: string;

  severityText: string;
  severityNumber: number;

  message: string;

  status: LogStatus;
};

export type LogStatus = {
  kind: LogStatusKind;
  reason?: string;
  objectID?: string;
  skipNotification: boolean;
  skipStorage: boolean;
};

export enum LogStatusKind {
  LogActive = 0,
  LogSnoozed,
  LogSilenced,
  LogRatelimited,
  LogDropped,
  LogActiveCheck,
  LogAcked,
}

export type GetAutocompleteLogData = {
};

export function getAutocompleteLog(data: GetAutocompleteLogData): Promise<string[]> {
  var q = {
    search: data.search,
    field: data.field,
  };
  return axios.get<string[]>("/api/autocomplete/log", { data: q}).then((resp) => {
    return resp.data;
  });
}

