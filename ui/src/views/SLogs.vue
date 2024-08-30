<script setup lang="ts">
import { h, ref, onMounted, onActivated, Component } from 'vue'
import axios from 'axios'
import type { AxiosResponse } from 'axios'
import { NIcon, NButton, NSpace, NDataTable, NInputGroup, NLoadingBarProvider, useLoadingBar } from 'naive-ui'
import { Refresh } from '@vicons/ionicons5'

import SSearch from '@/components/SSearch.vue'
import STimeRange from '@/components/STimeRange.vue'
import STimestamp from '@/components/STimestamp.vue'
import STimestampTitle from '@/components/STimestampTitle.vue'
import SLogAttributes from '@/components/SLogAttributes.vue'
import type { Log, LogResults } from '@/api/types'
import { defaultRangeMillis } from '@/utils/timerange'
import { usePagination } from '@/utils/pagination'

const search = ref<string>("")
const items = ref<Array<Log>>()
const rangeMillis = ref<[number, number]>(defaultRangeMillis())
const stimerange = ref(null)
const loading = useLoadingBar()
const pagination = usePagination(getLogs)
const table = ref<undefined|HTMLElement>(undefined)

function getLogs(): Promise {
  // loading.value = true
  loading.start()
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
        pagination.setMore(resp.data.more)
      } else {
        console.log("Logs not found")
      }
      // loading.value = false
      loading.finish()
    })
    .catch((err) => {
      // items.value = []
      // loading.value = false
      loading.error()
    })
}

function renderExpand(row) {
  console.log(`renderExpand: ${row}`)
  console.log(row)
  return h("pre", null, JSON.stringify(row, null, 2))
}

function render(component: Component, attr: string) {
  return (row) => {
    var options = {}
    options[attr] = row[attr]
    return h(component, options)
  }
}

const columns = [
  {type: 'expand', renderExpand: renderExpand},
  {
    key: 'timestampMillis',
    title: () => h(STimestampTitle),
    render: render(STimestamp, "timestampMillis"),
    width: 150,
  },
  {title: 'Attributes', render: (row) => h(SLogAttributes, {row: row}), width: 300},
  {title: 'Message', key: 'message', ellipsis: {tooltip: {placement: "bottom-end", width: 500}}},
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

function rowProps(row: Log) {
  return {
    onContextmenu: (e) =>{
      console.log("right click supported!")
      e.preventDefault()
      return false
    },
  }
}

</script>

<template>
  <div>
    <n-input-group>
      <s-time-range ref="stimerange" v-model:rangeMillis="rangeMillis" @update="onUpdateTimerange" />
      <s-search @search="onSearch" />
      <n-button @click="getLogs()"><n-icon :component="Refresh" /></n-button>
    </n-input-group>
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
  </div>
</template>
