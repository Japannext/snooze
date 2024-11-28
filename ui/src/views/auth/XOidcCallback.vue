<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { OidcClient, UserManager } from 'oidc-client-ts'
import { useRouteQuery } from '@vueuse/router'

import { getAuthMethods, type AuthMethod, useOidc } from '@/api'

// const method = ref<AuthMethod>(null)
const methodName = useRouteQuery("method")

onMounted(() => {
  console.log(`methodName: ${methodName.value}`)
  getAuthMethods()
    .then(methods => {
      for (var method of methods) {
        if (method.name == methodName.value) {
          return method
        }
      }
      throw new Error(`could not find method ${methodName.value}`)
    })
    .then(method => {
      useOidc(method).callback()
    })
})

</script>

<template>
  <div>Welcome</div>
</template>
