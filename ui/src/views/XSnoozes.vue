<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useLoadingBar, useMessage } from 'naive-ui'
import { usePagination } from '@/api'

// Components
import { NIcon, NButton, NButtonGroup, NSpace, NDataTable, NInputGroup, type DataTableColumn, type SelectOption, type RadioButtonProps } from 'naive-ui'
import { XSearch, XFilter, XTimeRange, XTimestampTitle } from '@/components/interface'
import { XTime, XDuration, XTagList, XSnoozeTime } from '@/components/attributes'
import { XNewSnoozeModal, XCancelSnoozeModal } from '@/components/modal'
import { Refresh, Add } from '@/icons'

import { getSnoozes, type Snooze, type GetSnoozesParams, type Filter } from '@/api'

const items = ref<Snooze[]>()
const selectedItems = ref<string[]>([])
const xTimerange = ref()

const loading = useLoadingBar()
const message = useMessage()
const showNewSnoozeModal = ref<boolean>(false)
const showCancelSnoozeModal = ref<boolean>(false)

const params = ref<GetSnoozesParams>({
  search: "",
  pagination: usePagination(refresh),
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
    key: 'time',
    title: 'Time constraint',
    render: renderSnoozeTime,
  },
  {
    key: 'duration',
    title: 'Duration',
    render: (row) => h(XDuration, {duration: (row.expireAt - row.startAt)}),
    width: 150,
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
  getSnoozes(params.value)
    .then((list) => {
      items.value = list.items
      params.value.pagination.itemCount = list.total
      params.value.pagination.setMore(list.more)
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
    start: row.startAt,
    end: row.expireAt,
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
  const tags = row.tags.map(x => {
    return {name: x, description: "", color: ""}
  })
  return h(XTagList, {tags: tags})
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
    <x-new-snooze-modal v-model:show="showNewSnoozeModal" @success="refreshWithDelay" />
    <x-cancel-snooze-modal v-model:show="showCancelSnoozeModal" :ids="selectedItems" @success="refreshWithDelay" />
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
