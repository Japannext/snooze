<script setup lang="ts">

import {
  h, ref, reactive, computed,
  useSlots, withDefaults, onMounted,
} from 'vue'
import type { Component } from 'vue'
import {
  NDataTable,
  NButton, NButtonGroup, NIcon,
  NGrid, NGi, NSpace,
} from 'naive-ui'
import { Refresh, Pencil } from '@vicons/ionicons5'
import type { DataTableRowKey, DataTableColumn, DataTableSortState } from 'naive-ui'

import { SSearch, SItemView, SButtonCreate, SModalEdit } from '@/components'
import type { Direction } from '@/api'
import type { SearchParams, CrudService } from '@/types'

const data = ref([])

interface Props {
  name: string,
  columns: DataTableColumn[],
  service: CrudService,
  form: Component,
  crud: boolean,
  renderView?,
}

const props = withDefaults(defineProps<Props>(), {
})

const slots = useSlots()

const extraActions = computed(() => slots.actions ? slots.actions() : undefined)
const modalEdit = ref(null)

function renderActions(row) {
  return h(NButtonGroup, [
    h(NButton, {
      type: "warning",
      round: true,
      size: "small",
      renderIcon: () => h(NIcon, {component: Pencil}),
      onClick: () => modalEdit.value.show(row.oid, row),
    })
  ])
}

function renderItemView(row) {
  return h(SItemView, {oid: row.oid, service: props.service})
}

const preppedColumns = [
  {type: 'selection'},
  {
    type: 'expand',
    renderExpand: props.renderView ?? renderItemView,
  },
  ...props.columns,
  {title: '', key: "snooze-web.actions", render: props.crud ? renderActions : undefined},
]

const formData = ref({})

// Parameters passed to the search
const params = ref<SearchParams>({
})

function refresh(): Promise<void> {
  return props.service.search(params.value)
  .then((collection) => {
    console.log("Refreshing...")
    collection.items.forEach(item => {
      console.log(JSON.stringify(item))
    })
    data.value = collection.items
  })
}

onMounted(() => {
  refresh({})
})

const showModalCreate = ref(false)
const showModalEdit = ref(false)

function submit() {
  // RuleService.ruleCreate(creatingRule.value)
}

const selection = ref<string[]>([])

function selectRows(rowKeys: DataTableRowKey[]) {
  selection.value = rowKeys
  console.log(`Selected ${JSON.stringify(selection.value)}`)
}

// Represent the visual cues given to the user concerning
// pagination
const pagination = reactive({
  page: 1,
  pageSize: 20,
  showSizePicker: true,
  pageSizes: [5, 20, 50],
})

/** Triggered when the query is changed.
 * @param {string} query The query to use for search
 */
function changeQuery(query: string) {
  params.value.query = query
  params.value.pageNb = 1
  refresh()
}

/** Change the page to a new page. Will first change the search
 *  parameters, trigger the search, then change the visual cue
 *  once the search is complete and the data reloaded.
 * @param {number} newPage The page to change to
 */
function changePage(newPage: number) {
  params.value.pageNb = newPage
  pagination.page = newPage
  refresh()
}

/** Change the page size. Will first change the search
 *  parameters, trigger the search, then change the visual cue
 *  once the search is complete and the data reloaded.
 *  Will also reset the page counter.
 * @param {number} newPageSize The new page size
 */
function changePageSize(newPageSize: number) {
  params.value.pageSize = newPageSize
  params.value.page = 1
  pagination.page = 1
  pagination.pageSize = newPageSize
  refresh()
}

/** Change the sorting column and the sorting direction.
 *  This function will change the search parameters, trigger
 *  the search, then change the sorting visual cue once the
 *  search is done.
 * @param {DataTableSortState} sorter The new sorting rule
 *  requested by the user.
 */
function changeSort(sorter: DataTableSortState) {
  params.value.orderBy = sorter.columnKey
  switch(sorter.order) {
    case 'ascend':
      params.value.order = Direction['1']
      break
    case 'descend': default:
      params.value.order = Direction['-1']
      break
  }
  refresh()
  .then(() => {
    var col = props.columns.find(c => c.name == sorter.columnKey)
    if (!sorter) {
      col.sortOrder = false
    } else {
      col.sortOrder = sorter.order
    }
  })
}

</script>

<template>
  <n-grid :cols="1" y-gap="10">
    <n-gi>
      <n-grid :cols="3" x-gap="20">
        <n-gi>
          <slot name="left" />
        </n-gi>
        <n-gi>
          <s-search @search="changeQuery" />
        </n-gi>
        <n-gi>
          <n-space justify="right">
            <s-button-create
              v-if="crud"
              :name="name"
              :service="service"
              :form="form"
              @refresh="refresh"
            />
            <n-button-group>
              <n-button @click="refresh"><n-icon :component="Refresh" /></n-button>
            </n-button-group>
          </n-space>
        </n-gi>
      </n-grid>
    </n-gi>

    <n-gi>
      <n-data-table
        striped
        bordered
        :single-line="false"
        :columns="preppedColumns"
        :data="data"
        :row-key="(row) => row.oid"
        :pagination="pagination"
        :resizable="true"
        @update:checked-row-keys="selectRows"
        @update:sorter="changeSort"
        @update:page="changePage"
        @update:page-size="changePageSize"
      />
      <s-modal-edit
        ref="modalEdit"
        :name="name"
        :service="service"
        :form="form"
      />
    </n-gi>
  </n-grid>
</template>

<style>
.n-data-table .n-data-table-td {
  padding-top: 0px;
  padding-bottom: 0px;
}
</style>
