<script setup lang="ts">
import { defineProps, defineEmits, defineModel, ref} from 'vue'
import { NAlert, NModal, NCard, NForm, NGrid, NFormItemGi, NSpace, NButton, NInput, NSelect } from 'naive-ui'
import { createAck, type Ack } from '@/api/ack'

const show = defineModel('show', {type: Boolean, default: false})
const item = ref<Ack>({})
const buttonLoading = ref<Boolean>(false)
const alerts = ref<Array<String>>([])

const emit = defineEmits(['success'])

const props = defineProps<{
  ids: Array<string>,
}>()

const tagOptions = [
  {label: 'Maintenance', value: 'maintenance', class: 'warning'},
  {label: 'Incident', value: 'incident', class: 'error'},
]

async function create() {
  buttonLoading.value = true
  item.value.logIDs = props.ids
  createAck(item.value)
    .then(() => {
      show.value = false
      buttonLoading.value = false
      emit('success')
    })
    .catch((err) => {
      buttonLoading.value = false
      alerts.value += err.message
    })
}

function cancel() {
  show.value = false
  item.value = {}
}

</script>

<template>
  <n-modal
    v-model:show="show"
    style="width: 800px"
  >
    <n-card title="Acknowledge">
      <template #header-extra>
        acknowledge {{ ids.length }} logs
      </template>
      <n-form size="small">
        <n-grid>
          <n-form-item-gi label="Tags" :span="12">
            <n-select v-model="item.tags" multiple :options="tagOptions" />
          </n-form-item-gi>
          <n-form-item-gi label="Reason" :span="24">
            <n-input v-model:value="item.reason" type="textarea" round />
          </n-form-item-gi>
        </n-grid>
      </n-form>
      <n-alert
        v-for="alert in alerts"
        :key="alert"
        title="Error"
        type="error"
      >
        {{ alert }}
      </n-alert>
      <template #action>
        <n-space justify="end">
          <n-button @click="cancel">Cancel</n-button>
          <n-button type="primary" :loading="buttonLoading" @click="create">Acknowledge</n-button>
        </n-space>
      </template>
    </n-card>
  </n-modal>
</template>
