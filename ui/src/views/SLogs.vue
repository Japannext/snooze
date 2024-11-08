<script setup lang="ts">
import axios from 'axios'
import { h, ref, onMounted, onActivated, Component } from 'vue'
import { useLoadingBar, useMessage } from 'naive-ui'
import { defaultRangeMillis } from '@/utils/timerange'
import { usePagination } from '@/utils/pagination'

// Components
import { NRadioGroup, NRadioButton, NInputGroupLabel, NIcon, NTag, NCard, NTabs, NTab, NButton, NSpace, NDataTable, NInputGroup, NLoadingBarProvider } from 'naive-ui'
import SSearch from '@/components/SSearch.vue'
import STimeRange from '@/components/STimeRange.vue'
import STimestamp from '@/components/STimestamp.vue'
import STimestampTitle from '@/components/STimestampTitle.vue'
import SLogAttributes from '@/components/SLogAttributes.vue'
import SStatus from '@/components/SStatus.vue'
import SFilter from '@/components/SFilter.vue'
import { Refresh } from '@vicons/ionicons5'

// Types
import type { AxiosResponse } from 'axios'
import type { Log, LogResults } from '@/api/types'

const search = ref<string>("")
const items = ref<Array<Log>>()
const rangeMillis = ref<[number, number]>(defaultRangeMillis())
const stimerange = ref(null)
const loading = useLoadingBar()
const pagination = usePagination(getItems)
const table = ref<undefined|HTMLElement>(undefined)
const message = useMessage()
const filter = ref<string>("")

function getItems(): Promise {
  // loading.value = true
  loading.start()
  var timerange = stimerange.value.getTime()
  var params = {
    page: pagination.page,
    size: pagination.pageSize,
  }
  if (search.value) {
    params.search = search.value
  }
  if (timerange[0] > 0) {
    params.start = timerange[0]
  }
  if (timerange[1] > 0) {
    params.end = timerange[1]
  }
  if (filter.value) {
    params.filter = filter.value
  }
  console.log(`getItems(${JSON.stringify(params)})`)
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
    .catch(err => {
      // items.value = []
      // loading.value = false
      message.error(`failed to load logs: ${err}`)
      loading.error()
    })
}

function renderExpand(row) {
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
  {
    type: 'expand',
    renderExpand: renderExpand,
  },
  {
    title: 'Status',
    render: (row) => h(SStatus, {kind: row.status.kind}),
    width: 90,
  },
  {
    key: 'timestamp',
    title: () => h(STimestampTitle),
    render: (row) => h(STimestamp, {ts: row.displayTime}),
    width: 150,
  },
  {title: 'Attributes', render: (row) => h(SLogAttributes, {row: row}), width: 300},
  {title: 'Message', key: 'message', ellipsis: {
    tooltip: {placement: "bottom-end", width: 500},
    lineClamp: 2,
  }},
]

onMounted(() => {
  getItems()
})

function onUpdateTimerange() {
  pagination.page = 1
  getItems()
}

function onSearch(text: string) {
  search.value = text
  getItems()
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

const filters = {
  
}
</script>

<template>
  <div>
    <n-space>
      <n-input-group>
        <n-input-group-label>Filters</n-input-group-label>
        <n-radio-group v-model:value="filter" v-on:change="getItems">
          <n-radio-button label="Active" value="active" />
          <n-radio-button label="Snoozed" value="snoozed" />
          <n-radio-button label="Acked" value="acked" />
          <n-radio-button label="Failed" value="failed" />
        </n-radio-group>
      </n-input-group>
      <n-input-group>
        <s-time-range ref="stimerange" v-model:rangeMillis="rangeMillis" @update="onUpdateTimerange" />
        <s-search @search="onSearch" />
        <n-button @click="getItems()"><n-icon :component="Refresh" /></n-button>
      </n-input-group>
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
      :row-key="(row) => row._id"
      :pagination="pagination"
    />
  </div>
</template>
