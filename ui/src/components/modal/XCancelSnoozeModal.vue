<script setup lang="ts">
import { ref, defineProps, defineModel, defineEmits } from 'vue'
import { useMessage } from 'naive-ui'
import { cancelSnooze } from '@/api'

import { NModal, NSpace, NButton, NCard, NForm, NFormItemGi, NInput, NGrid } from 'naive-ui'

const message = useMessage()

const props = defineProps<{
  ids: string[],
}>()

const emit = defineEmits(['success'])

const show = defineModel('show', {type: Boolean, default: false})
const loading = ref<Boolean>(false)
const reason = ref<string>("")

function submit() {
  loading.value = true
  cancelSnooze(props.ids, reason.value)
    .then((resp) => {
      emit('success')
    })
    .catch((err) => {
      message.error(`failed to snooze: ${err}`)
    })
    .finally(() => {
      loading.value = false
      cancel()
    })
}

function cancel() {
  reason.value = ""
  show.value = false
}
</script>

<template>
  <n-modal
    v-model:show="show"
    style="width: 800px"
  >
    <n-card title="Cancel snooze">
      <template #header-extra>
        Cancel/delete a snooze
      </template>
      <n-form size="small">
        <n-grid>
          <n-form-item-gi label="Reason" size="small" :span="24">
            <n-input v-model:value="reason" type="textarea" />
          </n-form-item-gi>
        </n-grid>
      </n-form>
      <template #footer>
        <pre>{{ {ids: ids, reason: reason} }}</pre>
      </template>
      <template #action>
        <n-space justify="end">
          <n-button @click="cancel">Cancel</n-button>
          <n-button type="error" :loading="loading" @click="submit">Cancel the snooze</n-button>
        </n-space>
      </template>
    </n-card>
  </n-modal>
</template>
