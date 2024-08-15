<script setup lang="ts">
import { ref, computed, defineEmits, withDefaults } from 'vue'
import { useAsyncState } from '@vueuse/core'
import { NCard, NModal, NSpace, NIcon, NIconWrapper, NButton } from 'naive-ui'
import { Add } from '@vicons/ionicons5'
import type { Component } from 'vue'

import { useDevMode } from '@/stores'
import { SHttpAlert, SErrorBox } from '@/components'
import { HttpAlert } from '@/utils/error'
import { ValidationStatus } from '@/utils/validation'

const props = defineProps<{
  name: string,
  service,
  form: Component,
}>()

const showing = ref<boolean>(false)

function show() {
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

const { devMode } = useDevMode()

/** Exit the modal
 */
function exit() {
  newData.value = {}
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

const asyncCreate = useAsyncState(
  () => props.service.create(newData.value), [],
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
        <n-icon-wrapper :size="24">
          <n-icon :size="22" :component="Add" />
        </n-icon-wrapper>
        Create new {{ name }}
      </template>
      <component
        :is="form"
        v-model:value="newData"
        :validation="validation"
      />
      <template #footer>
        <s-error-box ref="errorbox" />
        <pre v-if="devMode">{{ newData }}</pre>
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
            :loading="asyncCreate.isLoading.value"
            @click="asyncCreate.execute"
          >
            Create
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
