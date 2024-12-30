<script setup lang="ts">
import { ref, defineProps, defineModel } from 'vue'
import { useMessage } from 'naive-ui'
import { createTag, type Tag } from '@/api'

import { NModal, NSpace, NColorPicker, NButton, NCard, NForm, NFormItemGi, NInput, NGrid } from 'naive-ui'
import { XTag } from '@/components'


defineProps<{
  show: boolean;
}>()

const show = defineModel('show', {type: Boolean, default: false})
const loading = ref<boolean>(false)
const message = useMessage()
const item = ref<Tag>(defaultTag())

const swatches = [
  '#a8071a', // red
  '#fa541c', // ???
  '#fa8c16', // orange
  '#ad4e00', // cognac
  '#faad14', // golden
  '#613400', // brown
  '#d4b106', // yellow
  '#a0d911', // lime
  '#237804', // green
  '#13c2c2', // cyan
  '#003eb3', // blue
  '#722ed1', // purple
  '#c41d7f', // magenta
]

function defaultTag(): Tag {
  return {
    name: "",
    description: "",
    color: "#7BB5FDFF",
  }
}

function create() {
  loading.value = true
  createTag(item.value)
    .catch((err) => {
      message.error(`failed to tag: ${err}`)
    })
    .finally(() => {
      show.value = false
      loading.value = false
    })
}

function cancel() {
  show.value = false
  item.value = defaultTag()
}
</script>

<template>
  <n-modal
    v-model:show="show"
    style="width: 800px"
  >
    <n-card>
      <template #header>New tag</template>
      <n-form size="small">
        <n-grid :x-gap="20">
          <n-form-item-gi label="Name" size="small" :span="12">
            <n-input v-model:value="item.name" />
          </n-form-item-gi>
          <n-form-item-gi label="Color" size="small" :span="12">
            <n-color-picker
              v-model:value="item.color"
              :swatches="swatches"
              :modes="['hex']"
              show-preview
            />
          </n-form-item-gi>
          <n-form-item-gi label="Description" size="small" :span="24">
            <n-input v-model:value="item.description" type="textarea" />
          </n-form-item-gi>
        </n-grid>
      </n-form>
      <template #footer>
        Preview: <x-tag :name="item.name" :color="item.color" />
        <pre>{{ item }}</pre>
      </template>
      <template #action>
        <n-space justify="end">
          <n-button @click="cancel">Cancel</n-button>
          <n-button
            type="success"
            :loading="loading"
            :disabled="!item.name"
            @click="create"
          >
            Create
          </n-button>
        </n-space>
      </template>
    </n-card>
  </n-modal>
</template>
