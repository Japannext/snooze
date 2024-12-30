<script setup lang="ts">
import { ref, defineModel, computed } from 'vue'
import { getGroups, type Group, type GetGroupsParams } from '@/api'
import { NSelect, type SelectOption } from 'naive-ui'

const groupOptions = ref<SelectOption[]>([])
const loading = ref(false)

const groupMap = ref<Map<string, Group>>(new Map<string, Group>())

const props = defineProps<{
  groups: Group[],
}>()

const emit = defineEmits<{
  (e: "update:groups", value: Group[]): void
}>()

const params = ref<GetGroupsParams>({
  search: undefined,
})

async function onFocus() {
  await getGroupOptions(undefined)
}

async function onSearch(query: string) {
  await getGroupOptions(query)
}

function getGroupOptions(query?: string) {
  loading.value = true
  if (query) {
    params.value.search = `*${query}*`
  }
  getGroups(params.value)
    .then((list) => {
      groupOptions.value = list.items.map(group => {
        groupMap.value.set(group.hash, group)
        var labels = Object.keys(group.labels).map(key => {
          return `${group.labels[key]}`
        })
        return {
          label: `[${group.name}] ${labels}`,
          value: group.hash,
        }
      })
    })
    .finally(() => {
      loading.value = false
    })
}

const groupHashes = computed(() => {
  return props.groups.map((group) => {
    return group.hash
  })
})

function updateSelection(hashes: string[]) {
  const selectedGroups = hashes.map((hash) => {
    var gr = groupMap.value.get(hash)
    if (gr) {
      return gr
    }
    return {name: "", hash: hash, id: undefined, labels: {}}
  })
  emit('update:groups', selectedGroups)
}

</script>

<template>
  <n-select
    multiple
    remote
    clearable
    filterable
    :value="groupHashes"
    :loading="loading"
    :options="groupOptions"
    @focus="onFocus"
    @search="onSearch"
    @update:value="updateSelection"
  />
</template>
