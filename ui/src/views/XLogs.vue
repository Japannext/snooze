<script setup lang="ts">
import axios from 'axios'
import { h, ref, onMounted, computed, Component } from 'vue'
import { useLoadingBar, useMessage } from 'naive-ui'
import { defaultRangeMillis } from '@/utils/timerange'
import { usePagination } from '@/utils/pagination'

// Components
import { NRadioGroup, NRadioButton, NInputGroupLabel, NIcon, NTag, NCard, NTabs, NTab, NButton, NSpace, NDataTable, NInputGroup, NLoadingBarProvider } from 'naive-ui'
import { XSearch, XFilter, XTimeRange, XTimestampTitle } from '@/components/interface'
import { XStatus, XTime, XLogAttributes } from '@/components/attributes'
import { XAckModal } from '@/components/modal'
import { Refresh } from '@/icons'

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
const filter = ref<string>("active")
const selectedItems = ref<Array<Log>>([])
const showAckModal = ref<Boolean>(false)

function ack() {
  showAckModal.value = true
}

const selectedIDs = computed(() => {
  return selectedItems.value.map((e) => e.ID)
})

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
    type: 'selection',
    // options: ['all', 'none'],
  },
  {
    title: 'Status',
    render: (row) => h(XStatus, {kind: row.status.kind}),
    width: 90,
  },
  {
    key: 'timestamp',
    title: () => h(XTimestampTitle),
    render: (row) => h(XTime, {ts: row.displayTime}),
    width: 150,
  },
  {title: 'Attributes', render: (row) => h(XLogAttributes, {row: row}), width: 300},
  {title: 'Message', key: 'message', ellipsis: {
    tooltip: {placement: "bottom-end", width: 500},
    lineClamp: 2,
  }},
  {
    type: 'expand',
    renderExpand: renderExpand,
  },
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
    }
  }
}

const filters = [
  {label: "Active", value: "active"},
  {label: "Snoozed", value: "snoozed"},
  {label: "Acked", value: "acked"},
  {label: "All", value: "all"},
]

</script>

<template>
  <div>
    <n-space style="padding: 5px; margin-bottom: 10px;">
      <x-filter v-model:value="filter" :filters="filters" @change="getItems" />
      <n-input-group>
        <x-time-range ref="stimerange" v-model:rangeMillis="rangeMillis" @update="onUpdateTimerange" />
        <x-search @search="onSearch" />
        <n-button @click="getItems()"><n-icon :component="Refresh" /></n-button>
      </n-input-group>
      <n-input-group>
        <template v-if="selectedItems.length > 0">
          <n-button type="primary" @click="ack()">Ack ({{ selectedItems.length }})</n-button>
          <n-button type="warning" @click="ack()">Escalate ({{ selectedItems.length }})</n-button>
          <x-ack-modal v-model:show="showAckModal" :ids="selectedIDs" />
        </template>
      </n-input-group>
    </n-space>
    <n-data-table
      ref="table"
      v-model:checked-row-keys="selectedItems"
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
