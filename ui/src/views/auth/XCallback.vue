<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouteParams } from '@vueuse/router'
import { useRouter } from 'vue-router'
import { getAuthCallback } from '@/api'

const router = useRouter()
const provider = useRouteParams("provider")

function getProvider(): string {
  if (typeof provider.value === 'string') {
    return provider.value
  } else if (Array.isArray(provider.value)) {
    return provider.value[0]
  }
  else { // null
    throw new Error(`provider in route is null`)
  }
}

onMounted(() => {
  var queryString = window.location.search
  var query = new URLSearchParams(queryString)
  var code = query.get("code")
  if (code == null) {
    throw new Error(`no code found in query`)
  }
  var state = query.get("state")
  if (state == null) {
    throw new Error(`no state found in query`)
  }
  getAuthCallback(getProvider(), code, state)
    .then(() => {
      window.history.pushState({},"", "/")
      router.push({name: "logs", query: {}})
    })
})
</script>

<template>
  <div>Welcome</div>
</template>
