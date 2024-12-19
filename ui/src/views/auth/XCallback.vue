<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouteParams } from '@vueuse/router'
import { useRouter } from 'vue-router'
import { getAuthCallback } from '@/api'

const router = useRouter()
const provider = useRouteParams("provider")

onMounted(() => {
  var queryString = window.location.search
  var query = new URLSearchParams(queryString)
  var code = query.get("code")
  var state = query.get("state")
  getAuthCallback(provider.value, code, state)
    .then(() => {
      window.history.pushState({},"", "/")
      router.push({name: "logs", query: {}})
    })
})
</script>

<template>
  <div>Welcome</div>
</template>
