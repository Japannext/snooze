<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { OidcClient, UserManager } from 'oidc-client-ts'
import { useRouteQuery } from '@vueuse/router'

import { getAuthMethods, type AuthMethod } from '@/api'

// const method = ref<AuthMethod>(null)
const methodName = useRouteQuery("method")

/*
onMounted(() => {
  console.log(`welcome to the callback`)
  if (methodName.value === undefined) {
    console.error(`no method selected`)
    return
  }
  var method: AuthMethod = null
  getAuthMethods()
    .then((methods) => {
      methods.forEach((m) => {
        if (m.name == methodName.value) {
          method = m
          return
        }
      })
      if (method == null) {
        console.error(`no method found for method=${methodName.value}`)
        return
      }
      var settings = {
        authority: method.oidc.url,
        client_id: method.oidc.clientID,
        redirect_uri: `/api/oidc/redirect?method=${method.name}`,
        scope: method.oidc.scopes,
        response_type: "code",
        response_mode: "query",
        loadUserInfo: true,
        monitorSession: true,
        onSigninCallback: () => {
          window.history.replaceState({}, document.title, window.location.pathname)
        },
      }
      console.log(`oidc settings (callback): ${JSON.stringify(settings)}`)
      new UserManager(settings).signinCallback(`/api/oidc/redirect?method=${method.name}`).then((user) => {
        console.log("callback", user)
      })
      .catch((err) => {
        console.error(err)
      })
    })
})
*/

function select(methods: AuthMethod[], name: string): AuthMethod {
  for (var method of methods) {
    if (method.name == name) {
      return method
    }
  }
  throw new Error(`could not find method ${name}`)
}

/*
onMounted(() => {
  var methods = await getAuthMethods()
  var method = select(methods, methodName)
  oidc = useOidc(method)
})
*/

</script>

<template>
  <div>Welcome</div>
</template>
