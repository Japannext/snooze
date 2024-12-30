<script setup lang="ts">
import { h, ref, defineModel, computed, type VNodeChild } from 'vue'
import { getTags, type GetTagsParams, type Tag } from '@/api'

import { NSelect, type SelectOption } from 'naive-ui'
import { XTag, XTagOption } from '@/components'

const tagOptions = ref<SelectOption[]>([])
const loading = ref(false)

const tagMap = ref<Map<string, Tag>>(new Map<string, Tag>())

const props = defineProps<{
  tags: Tag[],
}>()

const emit = defineEmits<{
  (e: "update:tags", value: Tag[]): void
}>()

const params = ref<GetTagsParams>({
  search: undefined,
  pagination: {},
})

async function onFocus() {
  await getTagOptions(undefined)
}

async function onSearch(search: string) {
  await getTagOptions(search)
}

function getTagOptions(query?: string): Promise<void> {
  loading.value = true
  if (query) {
    params.value.search = `*${query}*`
  }
  return getTags(params.value)
    .then((list) => {
      tagOptions.value = list.items.map(tag => {
        tagMap.value.set(tag.name, tag)
        return {
          value: tag.name,
        }
      })
    })
    .finally(() => {
      loading.value = false
    })
}

function renderLabel(option: SelectOption): VNodeChild {
  if (typeof option.value != 'string') {
    throw new Error(`unexpected type of option '${option}'`)
  }
  var tag = tagMap.value.get(option.value)
  if (tag) {
    return h(XTagOption, {name: tag.name, color: tag.color, description: tag.description})
  }
  return h(XTagOption, {name: option.value})
}

function renderTag({option, handleClose}: {option: SelectOption, handleClose: () => void}): VNodeChild {
  if (typeof option.value != 'string') {
    throw new Error(`selected option has no value: ${option}`)
  }
  var tag = tagMap.value.get(option.value)
  if (tag) {
    return h(XTag, {name: tag.name, color: tag.color, closable: true, onClose: handleClose })
  }
  return h(XTag, {name: option.value, closable: true, onClose: handleClose})
}

const tagNames = computed(() => {
  return props.tags.map((tag) => {
    return tag.name
  })
})

function updateSelection(names: string[]) {
  const selectedTags = names.map((name) => {
    var tag = tagMap.value.get(name)
    if (tag) {
      return tag
    }
    return {name: name, color: "", description: ""}
  })
  emit('update:tags', selectedTags)
}

</script>

<template>
  <n-select
    multiple
    remote
    clearable
    filterable
    :value="tagNames"
    :loading="loading"
    :options="tagOptions"
    :render-label="renderLabel"
    :render-tag="renderTag"
    @focus="onFocus"
    @search="onSearch"
    @update:value="updateSelection"
  />
</template>
