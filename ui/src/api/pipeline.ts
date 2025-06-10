import axios from "axios";
import type { Schedule } from "@/api"

export type Grouping = {
  name: string
  if?: string
  groupBy: GroupBy[]
}

export type GroupBy = {
  kind: string
  key: string
}

export type Mapping = {
  name: string;
  map: {[key: string]: string};
}

export type Field = {
  name: string
  nestedFields: string[]
}

export type Transform = {
  name: string
  if?: string
  actions: Action[]
}

export type Action = {
  set?: Set
  unset?: Unset
  regex?: Regex
}

export type Set = {
  identity?: {[key: string]: string};
  labels?: {[key: string]: string};
  message?: string;
}

export type Unset = {
  identity?: {[key: string]: string};
  labels?: {[key: string]: string};
  allIdentity?: boolean
  allLabels?: boolean
}

export type Regex = {
  field: string
  match: string
}

export type Profile = {
  name: string
  switch: Kv
  Patterns: Pattern
}

export type Kv = {
  key: string
  value: string
}

export type Pattern = {
  name: string
  description: string
  regex: string
  actions: Action[]
  groupBy: {[key: string]: string};
  droppedLabels: string[]
  drop: boolean
  silence: boolean
}

export type Silence = {
  name: string
  description?: string
  if: string
  schedule: Schedule
  drop: boolean
}

export type Ratelimit = {
  name: string
  if?: string
  group: string
  countPerPeriod: number
  period: string
}

export type Notification = {
  name: string
  if?: string
  destinations: Destination[]
}

export type Destination = {
  queue: string;
  profile: string;
};

export function getPipelineGroupings(): Promise<Grouping[]> {
  return axios.get<Grouping[]>("/api/pipeline/groupings").then((resp) => {
    return resp.data;
  });
}

export function getPipelineMappings(): Promise<Mapping[]> {
  return axios.get<Mapping[]>("/api/pipeline/mappings").then((resp) => {
    return resp.data;
  });
}

export function getPipelineTransforms(): Promise<Transform[]> {
  return axios.get<Transform[]>("/api/pipeline/transforms").then((resp) => {
    return resp.data;
  });
}

export function getPipelineProfiles(): Promise<Profile[]> {
  return axios.get<Profile[]>("/api/pipeline/profiles").then((resp) => {
    return resp.data;
  });
}

export function getPipelineSilences(): Promise<Silence[]> {
  return axios.get<Silence[]>("/api/pipeline/silences").then((resp) => {
    return resp.data;
  });
}

export function getPipelineRatelimits(): Promise<Ratelimit[]> {
  return axios.get<Ratelimit[]>("/api/pipeline/ratelimits").then((resp) => {
    return resp.data;
  });
}

export function getPipelineNotifications(): Promise<Notification[]> {
  return axios.get<Notification[]>("/api/pipeline/notifications").then((resp) => {
    return resp.data;
  });
}

export function getPipelineDefaultDestination(): Promise<Destination[]> {
  return axios.get<Destination[]>("/api/pipeline/default-destinations").then((resp) => {
    return resp.data;
  });
}
