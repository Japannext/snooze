<script setup lang="ts">
import { h, ref, onMounted, Component } from 'vue'
import axios from 'axios'
import type { AxiosResponse } from 'axios'
import {NModal, NSpace, NTabs, NTab, NCard, NDataTable, NInputGroup, useLoadingBar, useModal, useMessage } from 'naive-ui'

import SAlertAttributes from '@/components/SAlertAttributes.vue'
import STimeRange from '@/components/STimeRange.vue'
import STimestamp from '@/components/STimestamp.vue'
import SIdentity from '@/components/SIdentity.vue'
import type { Alert, AlertResults } from '@/api/types'
import { usePagination } from '@/utils/pagination'

const items = ref<Array<Alert>>([])
const loading = useLoadingBar()
const stimerange = ref(null)
const pagination = usePagination(listAlerts)
const message = useMessage()
const selected = ref<Alert>(null)
const showDetails = ref<boolean>(false)

function listAlerts(): Promise {
  loading.start()
  var timerange = stimerange.value.getTime()
  var params = {
    page: pagination.page,
    size: pagination.pageSize,
  }
  if (timerange[0] > 0) {
    params.start = timerange[0]
  }
  if (timerange[1] > 0) {
    params.end = timerange[1]
  }
  console.log(`listAlerts()`)
  return axios.get<Alert>("/api/alerts", {params: params})
    .then((resp: AxiosResponse<AlertResults>) => {
      if (resp.data) {
        items.value = resp.data.items
        pagination.itemCount = resp.data.total
        pagination.setMore(resp.data.more)
      } else {
        console.log("alerts not found")
      }
      loading.finish()
    })
    .catch(err => {
      message.error(`failed to load alerts: ${err}`)
      loading.error()
    })
}

function onUpdateTimerange() {
  pagination.page = 1
  listAlerts()
}

onMounted(() => {
  listAlerts()
})

function select(item: Alert) {
  selected.value = item
  showDetails.value = true
}

function renderExpand(row) {
  return h("pre", null, JSON.stringify(row, null, 2))
}

const columns = [
  {type: 'expand', renderExpand: renderExpand},
  {
    title: "Since",
    render: (row) => h(STimestamp, {value: row.startsAt}),
    width: 150,
  },
  {
    title: 'Attributes',
    render: (row) => h(SAlertAttributes, {row: row}),
  },
  {
    title: 'Summary',
    key: 'summary',
    ellipsis: {tooltip: {placement: "bottom-end", width: 500}},
  },
]

</script>

<template>
  <n-space>
    <s-time-range ref="stimerange" @update="onUpdateTimerange" />
    <n-tabs type="line">
      <n-tab name="Active" />
      <n-tab name="History" />
    </n-tabs>
  </n-space>
  <n-data-table
    ref="table"
    remote
    striped
    bordered
    :resizable="true"
    size="small"
    :single-line="false"
    :columns="columns"
    :data="items"
    :row-key="(row) => row.id"
    :pagination="pagination"
  />
</template>
