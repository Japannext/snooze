<script setup lang="ts">
import { ref, defineEmits, withDefaults } from 'vue'
import { useAsyncState } from '@vueuse/core'
import { NCard, NModal, NSpace, NIcon, NIconWrapper, NButton } from 'naive-ui'
import { Pencil } from '@vicons/ionicons5'
import type { Component } from 'vue'

import { SErrorBox } from '@/components'
import { HttpAlert } from '@/utils/error'
import { ValidationStatus } from '@/utils/validation'

const props = withDefaults(defineProps<{
  name: string,
  service,
  form: Component,
}>(), {})

const showing = ref<boolean>(false)
const oid = ref<string|null>(null)

function show(id: string, data: object) {
  oid.value = id
  newData.value = data
  showing.value = true
}

defineExpose({ show })

const emit = defineEmits<{
  // A signal emitted to tell the bigger view to refresh its view, since new
  // data was added
  (e: "refresh"): void,
}>()

const errorbox = ref(null)

const newData = ref<object>({})

const isDev: boolean = Boolean(import.meta.env.MODE == 'development')

/** Exit the modal
 */
function exit() {
  newData.value = {}
  oid.value = null
  showing.value = false
}

const validation = ref<ValidationStatus>(new ValidationStatus())

// Create HttpAlert from ApiError response, then append
// it to the error box.
function handleError(resp) {
  console.log(`handleError: ${JSON.stringify(resp)}`)
  errorbox.value.add(HttpAlert.fromResponse(resp))
  validation.value.onResponse(resp)
  console.log(`validation: ${JSON.stringify(validation.value.feedback)}`)
}

const asyncUpdate = useAsyncState(
  () => props.service.update({oids: [oid.value], item: newData.value}), [],
  {
    onSuccess: onSubmit,
    onError: handleError,
    immediate: false,
  },
)

const asyncValidate = useAsyncState(
  () => props.service.validate(newData.value), [],
  {
    onError: handleError,
    immediate: false,
  }
)

function onSubmit() {
  exit()
  emit("refresh")
}

</script>

<template>
  <n-modal :show="showing">
    <n-card>
      <template #header>
        <n-icon-wrapper color="warning" :size="24">
          <n-icon :size="18" :component="Pencil" />
        </n-icon-wrapper>
        Editing {{ name }} '{{ oid }}'
      </template>
      <component
        :is="form"
        v-model:value="newData"
        :validation="validation"
      />
      <template #footer>
        <s-error-box ref="errorbox" />
        <pre v-if="isDev">{{ newData }}</pre>
      </template>
      <template #action>
        <n-space justify="end" size="small">
          <n-button @click="exit()">Cancel</n-button>
          <n-button
            secondary
            type="warning"
            :loading="asyncValidate.isLoading.value"
            @click="asyncValidate.execute"
          >
            Validate
          </n-button>
          <n-button
            type="success"
            :loading="asyncUpdate.isLoading.value"
            @click="asyncUpdate.execute"
          >
            Update
          </n-button>
        </n-space>
      </template>
    </n-card>
  </n-modal>
</template>

<style scoped>
.n-card {
  min-width: 600px;
  max-width: 1200px;
}
</style>
