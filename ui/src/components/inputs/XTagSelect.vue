<script setup lang="ts">
import { h, ref, defineModel } from 'vue'
import { getTags, type Tag, type GetTagsParams } from '@/api'

import { NSelect, NSpace, NEllipsis } from 'naive-ui'
import { XTag } from '@/components/attributes'
import { XTagOption } from '@/components/inputs'

import type { VNodeChild } from 'vue'
import type { SelectBaseOption, SelectOption } from 'naive-ui'

const tags = defineModel('tags')
const tagOptions = ref<SelectOption>([])
const loading = ref<Boolean>(false)

const params = ref<GetTagsParams>({
  search: null,
})

async function onFocus() {
  await getTagOptions(null)
}

async function onSearch(query: string) {
  await getTagOptions(query)
}

function getTagOptions(query?: string): Promise {
  loading.value = true
  if (query) {
    params.value.search = `*${query}*`
  }
  getTags(params.value)
    .then((list) => {
      tagOptions.value = list.items.map(tag => {
        return {
          value: tag,
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

function renderTag({option, handleClose}): VNodeChild {
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
