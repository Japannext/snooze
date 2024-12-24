<script setup lang="ts">
import { h, ref, defineModel, type VNodeChild } from 'vue'
import { getTags, type GetTagsParams } from '@/api'

import { NSelect, type SelectOption } from 'naive-ui'
import { XTag } from '@/components/attributes'
import { XTagOption } from '@/components/inputs'

const tags = defineModel('tags')
const tagOptions = ref<SelectOption[]>([])
const loading = ref(false)

const params = ref<GetTagsParams>({
  search: undefined,
  timerange: {},
  pagination: {},
})

async function onFocus() {
  await getTagOptions(undefined)
}

async function onSearch(query: string) {
  await getTagOptions(query)
}

function getTagOptions(query?: string): Promise<void> {
  loading.value = true
  if (query) {
    params.value.search = `*${query}*`
  }
  return getTags(params.value)
    .then((list) => {
      tagOptions.value = list.items.map(tag => {
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
  return h(XTagOption, {tag: option.value})
}

function renderTag({option, handleClose}: {option: SelectOption, handleClose: () => void}): VNodeChild {
  return h(XTag, {tag: option.value, closable: true, onClose: handleClose })
}

</script>

<template>
  <n-select
    v-model:value="tags"
    multiple
    remote
    clearable
    filterable
    :loading="loading"
    :options="tagOptions"
    :render-label="renderLabel"
    :render-tag="renderTag"
    @focus="onFocus"
    @search="onSearch"
  />
</template>
