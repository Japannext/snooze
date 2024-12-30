import axios from 'axios'

import { type Tag } from '@/api'

export type Ack = {
  time: number
  username: string
  reason: string
  logIDs: string[]
  tags: Tag[]
}

export function createAck(item: Ack): Promise<void> {
  return axios.post("/api/ack", item)
}
