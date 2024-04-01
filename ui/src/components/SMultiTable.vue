<script setup lang="ts">
import { ref, VNodeChild } from 'vue'
import { NTabs, NTab } from 'naive-ui'

import { CrudService } from '@/types'
import { STable } from '@/components'

interface Table {
  title: string
  service: CrudService
  columns: Array<object>
  form: object
  icon: VNodeChild
  crud: boolean
}

const props = defineProps<{
  tables: Array<Table>,
}>()

const current = ref(props.tables[0])

</script>

<template>
  <s-table
    name="current.title"
    :columns="current.columns"
    :service="current.service"
    :form="current.form"
    :crud="current.crud"
  >
    <template #left>
      <n-tabs default-value="tables[0].title">
        <n-tab v-for="table in tables" :key="table.label" name="table.label">
          <n-icon :component="table.icon" />
          {{ table.label }}
        </n-tab>
      </n-tabs>
    </template>
  </s-table>
</template>
