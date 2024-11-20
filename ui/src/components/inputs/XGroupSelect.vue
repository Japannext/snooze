<script setup lang="ts">
import { ref, defineModel } from 'vue'
import { getGroups, type Group, type GetGroupsParams } from '@/api'

import { NSelect } from 'naive-ui'

import type { SelectOption } from 'naive-ui'

const groups = defineModel('groups')
const groupOptions = ref<SelectOption>([])
const loading = ref<Boolean>(false)

const params = ref<GetGroupsParams>({
  search: null,
})

async function onFocus() {
  await getGroupOptions(null)
}

async function onSearch(query: string) {
  await getGroupOptions(query)
}

function getGroupOptions(query?: string): Promise {
  loading.value = true
  if (query) {
    params.value.search = `*${query}*`
  }
  getGroups(params.value)
    .then((list) => {
      groupOptions.value = list.items.map(group => {
        var labels = Object.keys(group.labels).map(key => {
          return `${group.labels[key]}`
        })
        return {
          label: `[${group.name}] ${labels}`,
          value: group,
        }
      })
    })
    .finally(() => {
      loading.value = false
    })
}

</script>

<template>
  <n-select
    v-model:value="groups"
    multiple
    remote
    clearable
    filterable
    :loading="loading"
    :options="groupOptions"
    @focus="onFocus"
    @search="onSearch"
  />
</template>
