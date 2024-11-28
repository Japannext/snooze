import axios from 'axios'
import { ref, reactive } from 'vue'
// import { useOidc } from '@/api'

export type AuthMethod = {
  name: string
  displayName: string
  kind: string
  icon: string
  color: string

  oidc?: OidcMethod
}

export type Profile = {
  username: string
  role: string
}

export interface Auth {
  login(): Promise<Profile>
  // logout(): Promise<void>
  profile(): Profile
}

export type OidcMethod = {
  url: string
  clientID: string
  redirectURL: string
  scopes: string[]
}

/*
const auth = reactive({
  name: ref<string>(),
  method: ref<AuthMethod>(),
  auth: ref<Auth>(),
  setMethod(method: AuthMethod) {
    this.name.value = method.name
    this.method.value = method
    switch(method.kind) {
    case "oidc":
      this.auth.value = useOidc(method)
    }
  }
})
*/

export function getAuthMethods(): Promise<AuthMethod[]> {
  return axios.get("/api/auth/methods")
    .then((resp => {
      return resp.data
    }))
}

export async function useAuth() {
  const methods = ref<AuthMethod[]>([])

  await getAuthMethods()
    .then((items) => {
      console.log(`getAuthMethods()`)
      methods.value = items
    })

  function select(name: string): AuthMethod|null {
    for (var method of methods.value) {
      if (method.name == name) {
        return method
      }
    }
    return null
  }

  console.log(`methods: ${JSON.stringify(methods.value)}`)

  const auth = reactive({
    methods: methods,
    getMethods() {
      return methods.value
    },
    withOidc(name: string) {
      var method = select(name)
      if (method == null) {
        throw new Error(`no method for '${name}'`)
      }
      return useOidc(method)
    }
  })
  return auth
}
