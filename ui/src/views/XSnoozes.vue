<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useLoadingBar, useMessage } from 'naive-ui'
import { usePagination } from '@/api'

// Components
import { NIcon, NButton, NButtonGroup, NSpace, NDataTable, NInputGroup, type DataTableColumn, type SelectOption, type RadioButtonProps } from 'naive-ui'
import { XSearch, XFilter, XTimeRange, XTimestampTitle, XTime, XDuration, XTagList, XSnoozeTime, XModalSnoozeCreate, XModalSnoozeCancel, XGroupList } from '@/components'
import { Refresh, Add } from '@/icons'

import { getSnoozes, type Snooze, type GetSnoozesParams, type Filter } from '@/api'

const items = ref<Snooze[]>()
const selectedItems = ref<string[]>([])
const xTimerange = ref()

const loading = useLoadingBar()
const message = useMessage()
const showNewSnoozeModal = ref<boolean>(false)
const showCancelSnoozeModal = ref<boolean>(false)

const pagination = usePagination(refresh)
const params = ref<GetSnoozesParams>({
  search: "",
  pagination: {},
  filter: "active",
  timerange: {},
})

const filters: Filter[] = [
  {label: "Active", value: "active"},
  {label: "Upcoming", value: "upcoming"},
  {label: "Expired", value: "expired"},
  {label: "Cancelled", value: "cancelled"},
  {label: "All", value: "all"},
]

const columns: DataTableColumn<Snooze>[] = [
  {type: 'selection'},
  {type: 'expand', renderExpand: renderExpand},
  {
    key: 'from',
    title: 'From',
    render: (row) => h(XTime, {ts: row.startsAt, format: 'absolute'}),
    width: 150,
  },
  {
    key: 'to',
    title: 'To',
    render: (row) => h(XTime, {ts: row.endsAt, format: 'absolute'}),
    width: 150,
  },
  {
    key: 'groups',
    title: 'Groups',
    render: (row) => h(XGroupList, {groups: row.groups}),
  },
  {
    key: 'duration',
    title: 'Duration',
    render: (row) => h(XDuration, {duration: (row.endsAt - row.startsAt)}),
    width: 150,
  },
  {
    key: 'username',
    title: 'Username',
  },
  {
    title: 'Tags',
    key: 'tags',
    render: renderTags,
  },
  {title: 'Reason', key: 'reason', ellipsis: {tooltip: {placement: "bottom-end", width: 500}}},
]

onMounted(() => {
  refresh()
})

function refresh() {
  console.log(`refresh()`)
  loading.start()
  params.value.timerange = xTimerange.value.getTime()
  params.value.pagination = {
    page: pagination.page,
    pageSize: pagination.pageSize,
  }
  getSnoozes(params.value)
    .then((list) => {
      items.value = list.items
      pagination.itemCount = list.total
      pagination.setMore(list.more)
      loading.finish()
    })
    .catch((err) => {
      items.value = []
      message.error(`failed to load snoozes: ${err}`)
      loading.error()
    })
}

const delayMillis = 300

function renderSnoozeTime(row: Snooze) {
  var isCancelled: boolean = false
  if (row.cancelled !== undefined || row.cancelled != null) {
    isCancelled = true
  }
  return h(XSnoozeTime, {
    start: row.startsAt,
    end: row.endsAt,
    cancelled: isCancelled,
  })
}

function refreshWithDelay() {
  setTimeout(refresh, delayMillis)
}

function reset() {
  params.value.pagination.page = 1
  refresh()
}

function renderTags(row: Snooze) {
  return h(XTagList, {tags: row.tags})
}

function renderExpand(row: Snooze) {
  return h("pre", null, JSON.stringify(row, null, 2))
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
      <n-button-group>
        <n-button type="warning" icon-placement="right" @click="showNewSnoozeModal = true">
          <template #icon><n-icon :component="Add" /></template>
          New Snooze
        </n-button>
        <n-button type="error" :disabled="selectedItems.length == 0" @click="showCancelSnoozeModal = true">
          Cancel ({{ selectedItems.length }})
        </n-button>
      </n-button-group>
    </n-space>
    <x-modal-snooze-create v-model:show="showNewSnoozeModal" @success="refreshWithDelay" />
    <x-modal-snooze-cancel v-model:show="showCancelSnoozeModal" :ids="selectedItems" @success="refreshWithDelay" />
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
      :pagination="pagination"
    />
  </div>
</template>
