<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouteParams, useRouteHash, useRouteQuery } from '@vueuse/router'
import { useRoute } from 'vue-router'
import { getAuthCallback } from '@/api'
import { xSnoozeToken } from '@/stores'

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
    .then(user => {
      console.log(user)
      xSnoozeToken.value = user.AccessToken
    })
})
</script>

<template>
  <div>Welcome</div>
</template>
