import axios from 'axios'
import type { AxiosResponse } from 'axios'

export type LoginBackend = {
  kind?: LoginBackendKind;
}

export enum LoginBackendKind {
  Oidc = "oidc",
  Static = "static",
}

export type LoginBackendList = {
  items: Array<LoginBackend>;
}

export function getLoginBackends(): Promise<LoginBackendList> {
  return axios.get<LoginBackendList>('/login-backends')
    .then(function(resp: AxiosResponse<LoginBackendList>) {
      return resp.data
    })
}
