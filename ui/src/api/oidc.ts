import { ref, unref, reactive } from 'vue'
import { createOidc } from 'oidc-spa'
import { AuthMethod, OidcMethod, getAuthMethods } from '@/api'
import { xSnoozeToken } from '@/stores'

export class OidcBackend {

  prOidc

  constructor(name: string, oidcMethod: OidcMethod) {
    var settings = {
      issuerUri: oidcMethod.url,
      clientId: oidcMethod.clientID,
      publicUrl: '/',
    }
    this.prOidc = createOidc(settings)
  }

  async login(): Promise<void> {
    var oidc = await this.prOidc
    if (!oidc.isUserLoggedIn) {
      return oidc.login({doesCurrentHrefRequiresAuth: false})
    } else {
      const { accessToken, decodedIdToken } = oidc.getTokens()
      xSnoozeToken.value = accessToken
    }
  }
}

export const oidcs = reactive({
  backends: new Map<string, OidcBackend>(),
  initialized: ref<boolean>(false),
  init(methods: AuthMethod[]): void {
    if (this.initialized) {
      return
    }
    for (var method of methods) {
      if (method.kind == "oidc" && method.oidc !== undefined) {
        var backend = new OidcBackend(method.name, method.oidc)
        this.backends.set(method.name, backend)
      }
    }
    this.initialized = true
  },
  get(name: string): OidcBackend|undefined {
    return this.backends.get(name)
  }
})

/*
export async function useOidc(method: AuthMethod): Promise<Auth> {

  if (method.oidc === undefined) {
    throw new Error(`unexpected auth config: 'oidc' field is undefined for method '${method.name}'`)
  }
  var settings = {
    issuerUri: method.oidc.url,
    clientId: method.oidc.clientID,
    publicUrl: '/',
  }
  const oidc = await createOidc(settings)

  const auth = reactive({
    login() {
      console.log("oidc-spa.login()")
      if (!oidc.isUserLoggedIn) {
        console.log("not logged in")
        return oidc.login({doesCurrentHrefRequiresAuth: false})
      }
    }
  })

  return auth
}
*/

/*
Log.setLogger(console)
Log.setLevel(Log.INFO)
const store = new WebStorageStateStore({store: window.localStorage})
*/

/*
export function useOidc(method: AuthMethod): Auth {
  var settings = {
    issuerUri: method.oidc.url,
    clientId: method.oidc.clientID,
    publicUrl: window.location.origin,
  }
  const auth = reactive({
    login(): Promise<void> {
      return createOidc(settings)
      .then(oidc => {
        oidc.login({doesCurrentHrefRequiresAuth: false})
      })
    }

  })
  const oidc = createOidc(settings)
}
*/

/*
export function useOidc(method: AuthMethod): Auth {
  if (method.oidc === undefined) {
    throw new Error(`unexpected auth config: 'oidc' field is undefined for method '${method.name}'`)
  }
  const token = ref<string>("")
  const baseURL = window.location.origin
  var settings = {
    authority: method.oidc.url,
    client_id: method.oidc.clientID,
    scope: method.oidc.scopes.join(' '),
    response_type: "code" as const,
    // response_mode: "query" as const,
    automaticSilentRenew: false,
    filterProtocolClaims: false,
    loadUserInfo: true,
    includeIdTokenInSilentRenew : false,
    redirect_uri: `${baseURL}/#/oidc/callback?method=${method.name}`,
    // post_logout_redirect_uri: `${baseURL}/#/login`,
    // popup_redirect_uri: `${baseURL}/oidc/callback?method=${method.name}`,
    silent_redirect_uri: `${baseURL}/#/oidc/silent-callback?method=${method.name}`,

    onSigninCallback: () => {
      window.history.replaceState({}, document.title, window.location.pathname)
    }
  }
  const mgr = new UserManager(settings)
  const auth = reactive({
    token: token,
    mgr: mgr,
    login(): Promise<void> {
      console.log("oidc.login()")
      return mgr.signinRedirect()
      .then(() => {
        console.log("login success")
      })
    },
    callback(): Promise<void> {
      console.log("oidc.callback()")
      return mgr.signinRedirectCallback()
      .then(user => {
        console.log("callback success", user)
      })
    },
  })

  return auth
}
*/

/*
export function useOidc(method: AuthMethod) {
  if (method.oidc === undefined) {
    throw new Error(`unexpected method ${method.name} invoking useOidc with no oidc config`)
  }
  var baseURL = window.location.origin
  var storage = new WebStorageStateStore({ store: window.localStorage })
  var settings = {
    authority: method.oidc.url,
    client_id: method.oidc.clientID,
    redirect_uri: `${baseURL}/#/oidc/callback?method=${method.name}`,
    scope: method.oidc.scopes.join(' '),
    response_type: "code" as const,
    response_mode: "fragment" as const,
    userStore: storage,
    onSigninCallback: () => {
      window.history.replaceState({}, document.title, window.location.pathname)
    }
  }
  const mgr = ref<UserManager>(new UserManager(settings))
  const oidc = reactive({
    mgr: new UserManager(settings),
    login(): Promise<User> {
      return mgr.value.signinPopup()
    }
  })


  return oidc
}

export interface OidcState {
  user?: User
  token?: string
}

export function useOidc2(method: AuthMethod): OidcState {
  const state = reactive({
    settings: method.oidc,
  })
}

async function autoAuthentication() {
}
*/
