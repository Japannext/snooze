<script setup lang="ts">
import { ref, defineProps, defineModel, defineEmits } from 'vue'
import { useMessage } from 'naive-ui'
import { createSnooze, type Snooze } from '@/api'

import { NModal, NSpace, NButton, NCard, NForm, NFormItemGi, NInput, NGrid } from 'naive-ui'
import { XTimePeriod, XGroupSelect, XTagSelect } from '@/components/inputs'

defineProps<{
  show: boolean;
}>()

const emit = defineEmits(['success'])

const show = defineModel('show', {type: Boolean, default: false})
const loading = ref<Boolean>(false)
const message = useMessage()
const item = ref<Snooze>(defaultValue())

function defaultValue(): Snooze {
  return {
    groups: [],
    tags: [],
  }
}

function submit() {
  loading.value = true
  console.log('submit()')
  createSnooze(item.value)
    .then((resp) => {
      emit('success')
      console.log('emit(`success`)')
    })
    .catch((err) => {
      console.log(`error: ${err}`)
      message.error(`failed to snooze: ${err}`)
    })
    .finally(() => {
      console.log('finally')
      loading.value = false
      exit()
    })
}

function exit() {
  show.value = false
  item.value = defaultValue()
}
</script>

<template>
  <n-modal
    v-model:show="show"
    style="width: 800px"
  >
    <n-card title="Add snooze">
      <n-form size="small">
        <n-grid>
          <n-form-item-gi label="Select" size="small" :span="24" required>
            <x-group-select v-model:groups="item.groups" />
          </n-form-item-gi>
          <n-form-item-gi label="Period" size="small" :span="24" required>
            <x-time-period v-model:start="item.startAt" v-model:end="item.expireAt" />
          </n-form-item-gi>
          <n-form-item-gi label="Reason" size="small" :span="24">
            <n-input v-model:value="item.reason" type="textarea" />
          </n-form-item-gi>
          <n-form-item-gi label="Tags" size="small" :span="24">
            <x-tag-select v-model:tags="item.tags" />
          </n-form-item-gi>
        </n-grid>
      </n-form>
      <template #footer>
        <pre>{{ item }}</pre>
      </template>
      <template #action>
        <n-space justify="end">
          <n-button @click="exit">Cancel</n-button>
          <n-button type="warning" :loading="loading" @click="submit">Snooze</n-button>
        </n-space>
      </template>
    </n-card>
  </n-modal>
</template>
