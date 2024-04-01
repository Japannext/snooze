<script setup lang="ts">
import { ref } from 'vue'
import { NLoadingBarProvider, useLoadingBar } from 'naive-ui'

import { AlertService, AlertV2, AlertGroupV2 } from '@/api'

const props = defineProps<{
  value: AlertGroupV2,
}>()

const loading = useLoadingBar()
const alerts = ref<AlertV2[]>([])
const query = ref<string>()

function refreshAlerts() {
  loading.start()
  AlertService.search()
  .then((collection) => {
    alerts.value = collection.items
    loading.finish()
  })
  .catch(() => {
    loading.error()
  })
}

const columns = props.value.ui_columns.map((key) => {
  return {
    title: key,
    key: key,
  }
})
</script>

<template>
  <n-loading-bar-provider>
    <pre>{{ value }}</pre>
    <n-data-table
      :columns="columns"
      :data="alerts"
    />
  </n-loading-bar-provider>
</template>
