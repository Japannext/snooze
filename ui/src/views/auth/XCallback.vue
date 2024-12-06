<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouteParams, useRouteHash, useRouteQuery } from '@vueuse/router'
import { useRouter } from 'vue-router'
import { getAuthCallback } from '@/api'

const router = useRouter()
const provider = useRouteParams("provider")

onMounted(() => {
  var queryString = window.location.search
  var query = new URLSearchParams(queryString)
  var code = query.get("code")
  var state = query.get("state")
  console.log(`query: ${JSON.stringify(query)}`)
  console.log(`code: '${code}'`)
  console.log(`state: '${state}'`)
  getAuthCallback(provider.value, code, state)
    .then(() => {
      router.push("logs")
    })
})
</script>

<template>
  <div>Welcome</div>
</template>
