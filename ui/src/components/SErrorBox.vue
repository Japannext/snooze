<script setup lang="ts">
// A box to report form validation errors, HTTP call errors, etc.
// It can be included in the form, or put on a side menu.
import { reactive } from 'vue'
import { NSpace, NAlert, NP } from 'naive-ui'

import { Alert } from '@/utils/error'

const alerts = reactive<Alert[]>([])

function add(newAlert: Alert) {
  console.log(`SErrorBox.add(${JSON.stringify(newAlert)})`)
  alerts.unshift(newAlert)
}

function close(idx: number) {
  console.log(`SErrorBox.close(${idx})`)
  alerts.splice(idx, 1)
}

defineExpose({ add })
</script>

<template>
  <n-space>
    <n-alert
      v-for="(alert, idx) in alerts"
      :key="alert.id"
      :type="alert.type"
      :title="alert.title"
      closable :on-close="() => close(idx)"
    >
      {{ alert.text }}
    </n-alert>
  </n-space>
</template>

<style scoped>

/* Ensure the multi-line text within the n-alert is displayed properly */
.n-alert {
  white-space: pre-line;
  word-wrap: break-word;
}

</style>
