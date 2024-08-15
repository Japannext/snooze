import axios from 'axios'
import type { AxiosResponse } from 'axios'

export type AuthButton = {
  displayName: string;
  path: string;
  icon?: string;
  color?: string;
}

export type AuthButtons = {
  buttons: Array<AuthButton>;
}

export function getAuthButtons(): Promise<AxiosResponse<AuthButtons>> {
  return axios.get<AuthButtons>('/api/auth/buttons')
    .then(function(resp: AxiosResponse<AuthButtons>) {
      return resp
    })
}
