import axios from 'axios'
import type { Destination, Schedule } from '@/api'

export type ProcessConfig = {
  transforms: Transform[]
  groupings: Grouping[]
  profiles: Profile[]
  silences: Silence[]
  ratelimits: Ratelimit[]
  notifications: Notification[]
  defaultDestinations: Destination[]
}

export type Transform = {
  name: string
  if?: string
}

export type Grouping = {
  name: string
  if?: string
  groupBy: string[]
}

export type Profile = {
  name: string
  switch: Kv
  patterns: Pattern[]
}

export type Kv = {
  key: string
  value: string
}

export type Pattern = {
  name: string
  regex: string
  groupBy: object
  identityOverride?: object
  droppedLabels?: string[]
  extraLabels?: object
  drop: boolean
  silence: boolean
  message: string
}

export type Silence = {
  name: string
  if: string
  schedule?: Schedule
  drop: boolean
}

export type Ratelimit = {
  name: string
  if?: string
  group: string
  burst: number
  period: string
}

type Notification = {
  name: string
  if?: string
  destinations: Destination[]
}

export function getProcessConfig(): Promise<ProcessConfig> {
  return axios.get('/process/config')
  .then((resp) => {
    return resp.data
  })
}
