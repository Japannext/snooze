import axios from 'axios'

export type Ack = {
  time: number;
  username: string;
  reason: string;
  logIDs: string[];
}

export function createAck(item: Ack): Promise<void> {
  return axios.post("/api/ack", item)
}
