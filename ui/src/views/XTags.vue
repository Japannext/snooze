<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useLoadingBar, useMessage } from 'naive-ui'
import { usePagination } from '@/api'

// Components
import { NIcon, NButton, NButtonGroup, NSpace, NDataTable, NInputGroup, type DataTableColumn } from 'naive-ui'
import { XSearch, XFilter, XTimeRange, XTag, XModalTagCreate } from '@/components'
import { Refresh, Add, Trash } from '@/icons'

import { getTags, type Tag, type GetTagsParams } from '@/api'

const items = ref<Tag[]>()
const selectedItems = ref<string[]>([])

const loading = useLoadingBar()
const message = useMessage()
const showNewTagModal = ref<boolean>(false)
const showDeleteTagModal = ref<boolean>(false)

const pagination = usePagination(refresh)

const params = ref<GetTagsParams>({
  search: "",
  pagination: {},
})

const columns: DataTableColumn<Tag>[] = [
  {type: 'selection'},
  {type: 'expand', renderExpand: renderExpand},
  {
    key: 'name',
    title: 'Name',
    render: (row) => h(XTag, {name: row.name, color: row.color}),
    width: 150,
  },
  {
    title: 'Description',
    key: 'description',
  },
  {
    title: 'Actions',
    key: 'actions',
  },
]

onMounted(() => {
  refresh()
})

function refresh() {
  loading.start()
  params.value.pagination = {
    page: pagination.page,
    pageSize: pagination.pageSize,
  }
  getTags(params.value)
    .then((list) => {
      items.value = list.items
      pagination.itemCount = list.total
      pagination.setMore(list.more)
      loading.finish()
    })
    .catch((err) => {
      items.value = []
      message.error(`failed to load tags: ${err}`)
      loading.error()
    })
}

function reset() {
  pagination.page = 1
  refresh()
}

function renderExpand(row: Tag) {
  return h("pre", null, JSON.stringify(row, null, 2))
}

</script>

<template>
  <div>
    <n-space :size="100" justify="center" style="padding: 5px; margin-bottom: 10px;">
      <n-input-group>
        <x-search v-model:value="params.search" @change="reset" />
        <n-button @click="refresh"><n-icon :component="Refresh" /></n-button>
      </n-input-group>
      <n-button-group>
        <n-button type="success" icon-placement="right" @click="showNewTagModal = true">
          <template #icon><n-icon :component="Add" /></template>
          New Tag
        </n-button>
        <n-button type="error" :disabled="selectedItems.length == 0" @click="showDeleteTagModal = true">
          <template #icon><n-icon :component="Trash" /></template>
          Delete ({{ selectedItems.length }})
        </n-button>
      </n-button-group>
    </n-space>
    <x-modal-tag-create v-model:show="showNewTagModal" />
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
