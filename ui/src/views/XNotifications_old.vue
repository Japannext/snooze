<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import axios from 'axios'
import type { AxiosResponse } from 'axios'
import { NTag, NIcon, NButton, NDataTable, NInputGroup } from 'naive-ui'
import { useLoadingBar } from 'naive-ui'
import { Refresh } from '@vicons/ionicons5'

import SSearch from '@/components/SSearch.vue'
import STimeRange from '@/components/STimeRange.vue'
import STimestamp from '@/components/STimestamp.vue'
import STimestampTitle from '@/components/STimestampTitle.vue'
import SDestination from '@/components/SDestination.vue'
import type { Notification, NotificationResults } from '@/api/types'
import { defaultRangeMillis } from '@/utils/timerange'
import { usePagination } from '@/utils/pagination'

const search = ref<string>("")
const items = ref<Array<Notification>>()
const loading = useLoadingBar()
const rangeMillis = ref<[number, number]>(defaultRangeMillis())
const stimerange = ref(null)
const pagination = usePagination(getItems)

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
  console.log(`getItems(${JSON.stringify(params)})`)
  return axios.get("/api/notifications", {params: params})
    .then((resp: AxiosResponse<NotificationResults>) => {
      if (resp.data) {
        items.value = resp.data.items
        pagination.itemCount = resp.data.total
        pagination.setMore(resp.data.more)
      } else {
        console.log("Notifications not found")
      }
      loading.finish()
    })
    .catch((err) => {
      // items.value = []
      loading.error()
    })
}

function renderExpand(row) {
  console.log(`renderExpand: ${row}`)
  console.log(row)
  return h("pre", null, JSON.stringify(row, null, 2))
}

const columns = [
  {type: 'expand', renderExpand: renderExpand},
  {
    key: 'timestamp.display',
    title: () => h(STimestampTitle),
    render: (row) => h(STimestamp, {timestamp: row.timestamp}),
    width: 150,
  },
  {title: 'Destination', render: (row) => h(SDestination, {destination: row.destination})},
  {title: 'Body', key: 'body', ellipsis: {tooltip: {placement: "bottom-end", width: 500}}},
  {title: 'Action', width: 200},
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

//    @update:page="onPageChange"
//    @update:page-size="onPageSizeChange"
</script>

<template>
  <n-input-group>
    <s-time-range ref="stimerange" v-model:rangeMillis="rangeMillis" @update="onUpdateTimerange" />
    <s-search @search="onSearch" />
    <n-button @click="getItems()"><n-icon :component="Refresh" /></n-button>
  </n-input-group>
  <n-data-table
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
