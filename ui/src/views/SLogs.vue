<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import axios from 'axios'
import type { AxiosResponse } from 'axios'
import { NTag, NIcon, NButton, NDataTable, NInputGroup } from 'naive-ui'
import { Refresh } from '@vicons/ionicons5'

import SSearch from '@/components/SSearch.vue'
import STimeRange from '@/components/STimeRange.vue'
import STimestamp from '@/components/STimestamp.vue'
import SIdentity from '@/components/SIdentity.vue'
import type { Log, LogResults } from '@/api/types'
import { defaultRangeMillis } from '@/utils/timerange'
import { usePagination } from '@/utils/pagination'

const search = ref<string>("")
const items = ref<Array<Log>>()
const loading = ref<boolean>(false)
const rangeMillis = ref<[number, number]>(defaultRangeMillis())
const stimerange = ref(null)
const pagination = usePagination(getLogs)

function getLogs(): Promise {
  loading.value = true
  var timerange = stimerange.value.getTime()
  var params = {
    page: pagination.page,
    size: pagination.pageSize,
  }
  if (search.value) {
    params.search = search
  }
  if (timerange[0] > 0) {
    params.start = timerange[0]
  }
  if (timerange[1] > 0) {
    params.end = timerange[1]
  }
  console.log(`getLogs(${JSON.stringify(params)})`)
  return axios.get<Log>("/api/logs", {params: params})
    .then((resp: AxiosResponse<LogResults>) => {
      if (resp.data) {
        items.value = resp.data.items
        pagination.itemCount = resp.data.total
      } else {
        console.log("Logs not found")
      }
      loading.value = false
    })
    .catch((err) => {
      items.value = []
      loading.value = false
    })
}

function renderSeverity(row) {
  var color: string
  // Trace
  if (row.severityNumber <= 4) {
    color = "default"
  }
  // Debug
  else if (row.severityNumber <= 8) {
    color = "default"
  }
  // Info
  else if (row.severityNumber <= 12) {
    color = "info"
  }
  // Warning
  else if (row.severityNumber <= 16) {
    color = "warning"
  }
  // Error
  else if (row.severityNumber <= 20) {
    color = "error"
  }
  // Fatal
  else {
    color = "error"
  }

  return h(NTag, {type: color}, {default: () => row.severityText})
}

function renderExpand(row) {
  console.log(`renderExpand: ${row}`)
  console.log(row)
  return h("pre", null, JSON.stringify(row, null, 2))
}

const columns = [
  {type: 'expand', renderExpand: renderExpand},
  {title: 'Timestamp', key: 'timestampMillis', render: (row) => h(STimestamp, {timestampMillis: row.timestampMillis}), width: 230},
  {title: 'Pattern', key: 'pattern', width: 200},
  {title: 'Severity', render: renderSeverity, width: 100, align: 'center'},
  {title: 'Identity', render: (row) => h(SIdentity, {identity: row.identity}), width: 300},
  {title: 'Message', key: 'message', ellipsis: {tooltip: {placement: "bottom-end", width: 500}}},
  {title: 'Action', width: 200},
]

onMounted(() => {
  getLogs()
})

function onUpdateTimerange() {
  pagination.page = 1
  getLogs()
}

function onSearch(text: string) {
  search.value = text
  getLogs()
}

//    @update:page="onPageChange"
//    @update:page-size="onPageSizeChange"
</script>

<template>
  <n-input-group>
    <s-time-range ref="stimerange" v-model:rangeMillis="rangeMillis" @update="onUpdateTimerange" />
    <s-search @search="onSearch" />
    <n-button @click="getLogs()"><n-icon :component="Refresh" /></n-button>
  </n-input-group>
  <n-data-table
    remote
    striped
    bordered
    :loading="loading"
    :resizable="true"
    size="small"
    :single-line="false"
    :columns="columns"
    :data="items"
    :row-key="(row) => row.id"
    :pagination="pagination"
  />
</template>
