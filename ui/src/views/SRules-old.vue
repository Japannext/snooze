<script setup lang="ts">

import { ref, reactive, onMounted } from 'vue'

import {
  NDataTable,
  NButton, NButtonGroup,
  NModal,
  NForm, NFormItem,
  NInput,
  NIcon, NIconWrapper,
  NGrid, NGi,
  NSpace,
} from 'naive-ui'
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'

import { Add, Refresh } from '@vicons/ionicons5'

import { RuleService } from '@/api'
import type { RuleV2, RuleV2Collection } from '@/api'

import SSearch from '@/components/SSearch.vue'

const columns = [
  {
    type: 'selection',
  },
  {title: "Name", key: "name"},
  {title: "Description", key: "description"},
]

const rules = ref([])

onMounted(() => {
  RuleService.ruleSearch()
    .then((collection: RuleV2Collection) => {
      rules.value = collection.items
    })
})

async function refresh() {
  RuleService.ruleSearch(
    pagination.page, pagination.pageSize,
  ).then((collection: RuleV2Collection) => {
    rules.value = collection.items
  })
}

const showCreate = ref(false)
const creatingRule = ref<RuleV2>({})

function submit() {
  RuleService.ruleCreate(creatingRule.value)
}

const selection = ref<string[]>([])

function selectRows(rowKeys: DataTableRowKey[]) {
  selection.value = rowKeys
}

const pagination = reactive({
  page: 0,
  pageSize: 20,
  showSizePicker: true,
  pageSizes: [5, 20, 50],
})

function changePage(currentPage: number) {
  refresh()
    .then(() => pagination.page = currentPage)
}

</script>

<template>
  <div>
    <n-grid :cols="1" y-gap="10">
      <n-gi>
        <n-grid :cols="3" x-gap="20">
          <n-gi />
          <n-gi>
            <s-search @search="search" />
          </n-gi>
          <n-gi>
            <n-space>
              <n-button-group>
                <n-button @click="refresh"><n-icon :component="Refresh" /></n-button>
              </n-button-group>
              <n-button-group>
                <n-button type="success" @click="showCreate = true"><n-icon :component="Add" />Add</n-button>
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
          :columns="columns"
          :data="rules"
          :row-key="(row) => row.name"
          :pagination="pagination"
          @update:checked-row-keys="selectRows"
          @update:sorter="sort"
          @update:page="changePage"
        />
      </n-gi>
    </n-grid>

    <n-modal
      v-model:show="showCreate"
      preset="dialog"
      positive-text="Create"
      negative-text="Cancel"
      @positive-click="submit"
      @negative-click="showCreate = false"
    >
      <template #icon>
        <n-icon-wrapper :size="24">
          <n-icon
            :size="22"
            :component="Add"
          />
        </n-icon-wrapper>
      </template>
      <template #header>
        Create rule
      </template>
      <n-form>
        <n-form-item
          label="Name"
          path="rule.name"
        >
          <n-input
            v-model:value="creatingRule.name"
            placeholder="Name"
          />
        </n-form-item>
        <n-form-item
          label="Description"
          path="rule.description"
        >
          <n-input
            v-model:value="creatingRule.description"
            placeholder="Name"
            type="textarea"
          />
        </n-form-item>
      </n-form>
      Creating rule: {{ creatingRule }}
    </n-modal>
  </div>
</template>
