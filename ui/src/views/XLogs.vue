<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useLoadingBar, useMessage } from 'naive-ui'
import { usePagination } from '@/api'

// Components
import { NIcon, NButton, NSpace, NDataTable, NInputGroup, type DataTableColumn, type SelectOption, RadioButtonProps } from 'naive-ui'
import { XSearch, XFilter, XTimeRange, XTimestampTitle } from '@/components/interface'
import { XStatus, XTime, XLogAttributes } from '@/components/attributes'
import { XAckModal } from '@/components/modal'
import { Refresh } from '@/icons'

import { getLogs, type Log, type GetLogsParams, type Filter } from '@/api'

const items = ref<Array<Log>>()
const selectedItems = ref<Array<string>>([])
const xTimerange = ref()

const loading = useLoadingBar()
const message = useMessage()
const showAckModal = ref<boolean>(false)
const showEscalateModal = ref<boolean>(false)

const params = ref<GetLogsParams>({
  search: "",
  filter: "active",
  pagination: usePagination(refresh),
  timerange: {},
})

const filters: Filter[] = [
  {label: "Active", value: "active"},
  {label: "Snoozed", value: "snoozed"},
  {label: "Acked", value: "acked"},
  {label: "All", value: "all"},
]

const columns: DataTableColumn<Log>[] = [
  {type: 'selection'},
  {
    type: 'expand',
    renderExpand: renderExpand,
  },
  {
    title: 'Status',
    render: (row) => h(XStatus, {status: row.status}),
    key: 'status',
    width: 90,
  },
  {
    key: 'timestamp',
    title: () => h(XTimestampTitle),
    render: (row) => h(XTime, {ts: row.displayTime}),
    width: 150,
  },
  {
    title: 'Attributes',
    render: (row) => h(XLogAttributes, {row: row}),
    width: 300,
    key: 'attributes',
  },
  {
    title: () => 'Message',
    key: 'message',
    ellipsis: {
      tooltip: {placement: "bottom-end", width: 500},
      lineClamp: 2,
    }
  },
]

onMounted(() => {
  refresh()
})

function refresh() {
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
      message.error(`failed to load logs (status ${err.response.status}): ${err.response.data}`, {duration: 10000})
      loading.error()
    })
}

function unselect() {
  selectedItems.value = []
  refresh()
}

function renderExpand(row: Log) {
  return h("pre", null, JSON.stringify(row, null, 2))
}

function reset() {
  params.value.pagination.page = 1
  refresh()
}

</script>

<template>
  <div>
    <n-space :size="100" justify="start" style="padding: 5px; margin-bottom: 10px;">
      <n-input-group>
        <x-time-range ref="xTimerange" @change="reset" />
        <x-search v-model:value="params.search" @change="reset" />
        <n-button @click="refresh"><n-icon :component="Refresh" /></n-button>
      </n-input-group>
      <x-filter v-model:value="params.filter" :filters="filters" @change="reset" />
      <n-input-group>
        <n-button type="primary" :disabled="selectedItems.length == 0" @click="showAckModal = true">Ack ({{ selectedItems.length }})</n-button>
        <x-ack-modal v-model:show="showAckModal" :ids="selectedItems" @success="unselect" />

        <n-button type="warning" :disabled="selectedItems.length == 0">Escalate ({{ selectedItems.length }})</n-button>
      </n-input-group>
    </n-space>
    <n-data-table
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

<style scoped>
th.n-data-table-th {
  padding: 0px;
}
</style>
