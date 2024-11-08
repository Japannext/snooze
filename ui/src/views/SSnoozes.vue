<script setup lang="ts">
import { h, ref, onMounted, onActivated, Component } from 'vue'
import axios from 'axios'
import type { AxiosResponse } from 'axios'
import { NIcon, NTag, NCard, NTabs, NTab, NButton, NSpace, NDataTable, NInputGroup, NLoadingBarProvider, useLoadingBar, useMessage, useModal } from 'naive-ui'
import { Refresh, Add } from '@vicons/ionicons5'
import { DateTime } from 'luxon'

import SSearch from '@/components/SSearch.vue'
import STimeRange from '@/components/STimeRange.vue'
import STimestamp from '@/components/STimestamp.vue'
import STimestampTitle from '@/components/STimestampTitle.vue'
import SLogAttributes from '@/components/SLogAttributes.vue'
import SSnoozeCreateModal from '@/components/modal/SSnoozeCreateModal.vue'
import type { Snooze, ListOf } from '@/api/types'
import { defaultRangeMillis } from '@/utils/timerange'
import { usePagination } from '@/utils/pagination'

const search = ref<string>("")
const items = ref<Array<Snooze>>()
const rangeMillis = ref<[number, number]>(defaultRangeMillis())
const stimerange = ref(null)
const loading = useLoadingBar()
const pagination = usePagination(getItems)
const table = ref<undefined|HTMLElement>(undefined)
const message = useMessage()

const showCreateModal = ref<Boolean>(false)

function getItems(): Promise {
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
  return axios.get<Snooze>("/api/snoozes", {params: params})
    .then((resp: AxiosResponse) => {
      // error handling
      return resp.data
    })
    .then((res: ListOf<object>) => {
      items.value = res.items
      pagination.itemCount = res.total
      pagination.setMore(res.more)
      loading.finish()
    })
    .catch(err => {
      message.error(`failed to load snoozes: ${err}`)
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
  {type: 'expand', renderExpand: renderExpand},
  {
    key: 'timestamp',
    title: () => h(STimestampTitle),
    render: (row) => h(STimestamp, {ts: row.startAt}),
    width: 150,
  },
  {title: 'Attributes', render: (row) => h(SLogAttributes, {row: row}), width: 300},
  {title: 'Message', key: 'message', ellipsis: {tooltip: {placement: "bottom-end", width: 500}}},
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

</script>

<template>
  <div>
    <n-input-group>
      <s-time-range ref="stimerange" v-model:rangeMillis="rangeMillis" @update="onUpdateTimerange" />
      <s-search @search="onSearch" />
      <n-button @click="showCreateModal = true"><n-icon :component="Add" /></n-button>
      <n-button @click="getItems()"><n-icon :component="Refresh" /></n-button>
    </n-input-group>
    <s-snooze-create-modal v-model:show="showCreateModal" />
    <n-tabs type="bar" tabs-padding="20" :theme-overrides="tabsStyleOverrides">
      <n-tab name="Active" />
      <n-tab name="Snoozed" />
      <n-tab name="Acknowledged" />
      <n-tab name="Failed">
        <n-space size="small">
          <span>Failed</span>
          <n-tag size="tiny" type="error" round>10</n-tag>
        </n-space>
      </n-tab>
    </n-tabs>
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
