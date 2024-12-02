import axios from 'axios'
export type AuthConfig = {
  genericOidc?: OidcConfig
}

export type OidcConfig = {
  clientID: string
  redirectURL: string
  scopes: string[]
  displayName: string
  icon: string
  color: string
}

export function getAuthConfig(): Promise<AuthConfig> {
  return axios.get("/api/auth/config")
    .then(resp => {
      return resp.data
    })
}

export function getAuthLogin(provider: string): Promise<void> {
  console.log(`getAuthLogin(${provider})`)
  return axios.get(`/api/auth/${provider}/login`)
    .then(resp => {
      console.log(resp.data)
      return
    })
}

export function getAuthCallback(provider: string, code: string, state: string): Promise<any> {
  console.log(`getAuthCallback(${provider})`)
  var params = {
    code: code,
    state: state,
  }
  return axios.get(`/api/auth/${provider}/callback`, {params: params})
    .then(resp => {
      return resp.data
    })
}
