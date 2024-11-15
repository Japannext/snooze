<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useLoadingBar, useMessage } from 'naive-ui'
import { usePagination } from '@/api'

// Components
import { NIcon, NButton, NSpace, NDataTable, NInputGroup } from 'naive-ui'
import { XSearch, XFilter, XTimeRange, XTimestampTitle } from '@/components/interface'
import { XStatus, XTime, XLogAttributes } from '@/components/attributes'
import { XAckModal } from '@/components/modal'
import { Refresh } from '@/icons'

import { getLogs, type Log, type GetLogsParams } from '@/api'

const items = ref<Array<Log>>()
const xTimerange = ref(null)
const loading = useLoadingBar()
const table = ref<undefined|HTMLElement>(undefined)
const message = useMessage()
const selectedItems = ref<Array<string>>([])

const showAckModal = ref<boolean>(false)
const showEscalateModal = ref<boolean>(false)

const params = ref<GetLogsParams>({
  search: "",
  filter: "active",
  pagination: usePagination(refresh)
})

const filters = [
  {label: "Active", value: "active"},
  {label: "Snoozed", value: "snoozed"},
  {label: "Acked", value: "acked"},
  {label: "All", value: "all"},
]

const columns = [
  {type: 'selection'},
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
  {
    title: () => 'Message',
    key: 'message',
    ellipsis: {
      tooltip: {placement: "bottom-end", width: 500},
      lineClamp: 2,
    }
  },
  {
    type: 'expand',
    renderExpand: renderExpand,
  },
]

function refresh(): Promise {
  loading.start()
  params.value.timerange = xTimerange.value.getTime()
  getLogs(params.value)
    .then((list) => {
      items.value = list.items
      params.value.pagination.itemCount = list.total
      params.value.pagination.setMore(list.more)
      loading.finish()
    })
    .catch((err) => {
      items.value = []
      message.error(`failed to load logs: ${err}`)
      loading.error()
    })
}

function unselect() {
  selectedItems.value = []
}

function renderExpand(row) {
  return h("pre", null, JSON.stringify(row, null, 2))
}

onMounted(() => {
  refresh()
})

function reset() {
  params.value.pagination.page = 1
  refresh()
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

</script>

<template>
  <div>
    <n-space :size="100" justify="center" style="padding: 5px; margin-bottom: 10px;">
      <x-filter v-model:value="params.filter" :filters="filters" @change="reset" />
      <n-input-group>
        <x-time-range ref="xTimerange" @change="reset" />
        <x-search v-model:value="params.search" @change="reset" />
        <n-button @click="refresh"><n-icon :component="Refresh" /></n-button>
      </n-input-group>
      <n-input-group>
        <n-button type="primary" :disabled="selectedItems.length == 0" @click="showAckModal = true">Ack ({{ selectedItems.length }})</n-button>
        <x-ack-modal v-model:show="showAckModal" :ids="selectedItems" @success="unselect" />

        <n-button type="warning" :disabled="selectedItems.length == 0">Escalate ({{ selectedItems.length }})</n-button>
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
      :pagination="params.pagination"
    />
  </div>
</template>

<style>
th.n-data-table-th {
  padding: 0px;
  background: red;
}
</style>
