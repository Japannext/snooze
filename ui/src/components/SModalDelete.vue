<script setup lang="ts">

import { ref, computed, withDefaults } from 'vue'
import {
  NCard,
  NModal, NSpace,
  NIcon, NIconWrapper, NButton,
} from 'naive-ui'
import { Add } from '@vicons/ionicons5'
import type { Component } from 'vue'

import { SHttpAlert } from '@/components'

interface Props {
  name: string,
  service,
  form: Component,
  show: boolean,
}

const props = withDefaults(defineProps<Props>(), {
  show: false,
})

const emit = defineEmits<{
  (e: "update:show", value: string): void,
}>()

const showValue = computed<boolean>({
  get() { return props.show },
  set(v) { emit("update:show", v) },
})

const newData = ref<object>({})

const isDev: boolean = Boolean(import.meta.env.MODE == 'development')

const validationStatus = ref<object>({})
const feedback = ref<object>({})
const submitLoading = ref<boolean>(false)
const validateLoading = ref<boolean>(false)

const showAlert = ref<boolean>(false)
const modalAlert = ref<object>({})

function resetLoading() {
  validateLoading.value = false
  submitLoading.value = false
}

function resetValidation() {
  showAlert.value = false
  validationStatus.value = {}
  feedback.value = {}
}

function handleError(resp) {
  resetLoading()
  modalAlert.value = {
    code: resp.status,
    name: resp.name,
    description: resp.body.detail,
  }
  if (resp.status == 400) {
    resetValidation()
    for (const [key, values] of Object.entries(resp.body.errors)) {
      const newKey = key.replace(/[.]__root__$/, '')
      validationStatus.value[newKey] = "error"
      feedback.value[newKey] = values.join("\n")

    }
  }
  showAlert.value = true
}

/** Exit the modal
 */
function exit() {
  showValue.value = false
  resetLoading()
  resetValidation()
  newData.value = {}
}

function create() {
  submitLoading.value = true
  props.service.create(newData.value)
  .then((_result) => {
    exit()
  }).catch(handleError)
}

function update() {
  submitLoading.value = true
  props.service.update(newData.value.oid, newData.value)
  .then((_result) => {
    exit()
  }).catch(handleError)
}

function validate() {
  validateLoading.value = true
  props.service.validate(newData.value)
  .then((_result) => {
    resetLoading()
  })
  .catch(handleError)
}

</script>

<template>
  <n-modal :show="showValue">
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
        :feedback="feedback"
        :validation-status="validationStatus"
      />
      <template #footer>
        <s-http-alert :show="showAlert" :alert="modalAlert" @update:show="resetValidation" />
        <pre v-if="isDev">{{ newData }}</pre>
      </template>
      <template #action>
        <n-space justify="end" size="small">
          <n-button @click="exit()">Cancel</n-button>
          <n-button secondary type="warning" @click="validate">Validate</n-button>
          <n-button type="success" @click="create">Create</n-button>
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
