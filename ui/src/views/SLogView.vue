<script setup lang="ts">
import axios from 'axios'
import type { AxiosResponse } from 'axios'
import { ref, onMounted } from 'vue'
import type { Log } from '@/api/types'
import { Refresh } from '@vicons/ionicons5'

const props = defineProps<{
  uid: string,
}>()

const item = ref<Log>

function fetchLog() {
  axios.get<Log>(`/api/log/${props.uid.value}`)
    .then((resp: AxiosResponse<Log>) => {
      if (resp.data) {
        item.value = resp.data
      } else {
        console.log("Log not found")
      }
    })
}

onMounted(() => {
  fetchLog()
})
</script>

<template>
  <n-card>
    <n-button @click="fetchLog()">
      <n-icon :component="Refresh" />
    </n-button>
    Timestamp: {{ item.timestamp }}
    Severity: {{ item.severity }}
  </n-card>
</template>
