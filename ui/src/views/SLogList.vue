<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import type { AxiosResponse } from 'axios'
import { NIcon, NButton, NCard, NDataTable, NSwitch } from 'naive-ui'
import { Refresh } from '@vicons/ionicons5'

import STimeRange from '@/components/STimeRange.vue'
import type { Log, LogsResponse } from '@/api/types'
import { defaultRangeMillis } from '@/utils/timerange'
import { usePagination } from '@/utils/pagination'

import { DateTime } from 'luxon'

const query = ref<string>("*")
// const timerange = ref<TimeRange>({start: 0, end: 0})
const items = ref<Array<Log>>()
const rangeMillis = ref<[number, number]>(defaultRangeMillis())

function fetchLogs() {
  var params = {
    page_nb: pagination.page,
    page_size: pagination.pageSize,
  }
  if (query.value) {
    params.query = query.value
  }
  if (rangeMillis.value[0] != 0) {
    params.start_millis = rangeMillis.value[0]
  }
  if (rangeMillis.value[1] != 0) {
    params.end_millis = rangeMillis.value[1]
  }
  axios.get<Log>("/api/logs", null, {params: params})
    .then((resp: AxiosResponse<LogsResponse>) => {
      if (resp.data) {
        items.value = resp.data.logs
      } else {
        console.log("Logs not found")
      }
    })
}

const timestampRelative = ref<boolean>(false)

function renderTimestamp(i: number): string {
  if (timestampRelative.value) {
    var diff = DateTime.fromMillis(i).diff(DateTime.now())
    return diff.rescale().toHuman()
  } else {
    return DateTime.fromMillis(i).toISO()
  }
}

function renderIdentity(row) {
  var identity = row.identity
  if (identity === undefined) {
    return "unknown"
  }
  switch (identity.kind) {
    case "host": {
      return `${identity.hostname} #${identity.process}`
    }
    case "k8s.pod": {
      return `${identity.pod} @${identity.namespace}`
    }
  }
}

const columns = [
  {title: 'Timestamp', key: 'timestampMillis', render: (row) => {return renderTimestamp(row.timestampMillis) }},
  {title: 'Log Pattern', key: 'logPattern'},
  {title: 'Severity', key: 'severityText'},
  {title: 'Identity', render: renderIdentity},
  {title: 'Message', key: 'body.message'},
]

onMounted(() => {
  fetchLogs()
})

// const pagination = ref<object>({page_nb: 1, page_size: 10, order_by: "timestampMillis", asc: false})

const pagination = usePagination()

/*
function updateSort(options: DataTableSortState | DataTableSortState[] | null) {
  if (options == null) {
    return
  }
  if (typeof options == DataTableSortState) {
    pagination.value.order_by = options.columnKey
    pagination.value.asc = (options.order == 'ascend')
    fetchLogs()
  }
}

function updatePageSize(pageSize: number) {
  pagination.value.page_size = pageSize
  fetchLogs()
}

function updatePage(page: number) {
  pagination.value.page_nb = page
  fetchLogs()
}
*/

</script>

<template>
  <n-button @click="fetchLogs()">
    <n-icon :component="Refresh" />
  </n-button>
  <n-switch v-model:value="timestampRelative">
    <template #checked>
      Relative time
    </template>
    <template #unchecked>
      ISO time
    </template>
  </n-switch>
  <s-time-range v-model:rangeMillis="rangeMillis" @refresh="fetchLogs()" />
  <n-data-table
    striped
    bordered
    :single-line="false"
    :columns="columns"
    :data="items"
    :row-key="(row) => row._id"
    :pagination="pagination"
    :resizable="true"
  />
  <n-card>
    <pre>{{ items }}</pre>
  </n-card>
</template>
