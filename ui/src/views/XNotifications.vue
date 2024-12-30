<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useLoadingBar, useMessage } from 'naive-ui'
import { usePagination } from '@/api'

// Components
import { NButton, NSpace, NIcon, NDataTable, NInputGroup, type DataTableColumn, type SelectOption, type RadioButtonProps } from 'naive-ui'
import { XSearch, XFilter, XTimeRange, XTimestampTitle, XDestination, XTime, XModalAck } from '@/components'
import { Refresh } from '@/icons'

import { getNotifications, type GetNotificationsParams, type Notification, type Filter } from '@/api'

const selected = ref<Notification|null>(null)
const showDetails = ref<boolean>(false)

// Utils
const loading = useLoadingBar()
const message = useMessage()

const items = ref<Array<Notification>>([])
const selectedItems = ref<Array<string>>([])
const xTimerange = ref()

const showAckModal = ref<boolean>(false)
const showEscalateModal = ref<boolean>(false)

const pagination = usePagination(refresh)
const params = ref<GetNotificationsParams>({
  search: "",
  filter: "active",
  pagination: {},
  timerange: {},
})

const filters: Filter[] = [
  {label: "Active", value: "active"},
  {label: "History", value: "history"},
]

const columns: DataTableColumn<Notification>[] = [
  {type: 'selection'},
  {
    title: () => h(XTimestampTitle),
    render: (row) => h(XTime, {ts: row.notificationTime}),
    key: 'timestamp',
    width: 150,
  },
  {
    title: 'Destination',
    render: (row) => h(XDestination, {destination: row.destination}),
    key: 'destination',
  },
  {type: 'expand', renderExpand: renderExpand},
]

function refresh() {
  loading.start()
  params.value.timerange = xTimerange.value.getTime()
  params.value.pagination = {
    page: pagination.page,
    pageSize: pagination.pageSize,
  }
  getNotifications(params.value)
    .then((list) => {
      items.value = list.items
      pagination.itemCount = list.total
      pagination.setMore(list.more)
      loading.finish()
    })
    .catch((err) => {
      items.value = []
      message.error(`failed to load alerts: ${err}`)
      loading.error()
    })
}

onMounted(() => {
  refresh()
})

function unselect() {
  selectedItems.value = []
}

function select(item: Notification) {
  selected.value = item
  showDetails.value = true
}

function reset() {
  params.value.pagination.page = 1
  refresh()
}

function renderExpand(row: Notification) {
  return h("pre", null, JSON.stringify(row, null, 2))
}

</script>

<template>
  <n-space :size="100" justify="center" style="padding: 5px; margin-bottom: 10px;">
    <x-filter v-model:value="params.filter" :filters="filters" @change="reset" />
    <n-input-group>
      <x-time-range ref="xTimerange" @change="reset" />
      <x-search v-model:value="params.search" @change="reset" />
      <n-button @click="refresh"><n-icon :component="Refresh" /></n-button>
    </n-input-group>
    <n-input-group>
      <n-button type="primary" :disabled="selectedItems.length == 0" @click="showAckModal = true">Ack ({{ selectedItems.length }})</n-button>
      <x-modal-ack v-model:show="showAckModal" :ids="selectedItems" @success="unselect" />

      <n-button type="warning" :disabled="selectedItems.length == 0">Escalate ({{ selectedItems.length }})</n-button>
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
    :pagination="params.pagination"
  />
</template>
