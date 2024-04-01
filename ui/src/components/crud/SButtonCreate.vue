<script setup lang="ts">
import { ref, defineEmits } from 'vue'
import { NButton, NIcon } from 'naive-ui'
import { Add } from '@vicons/ionicons5'

import { SModalCreate } from '@/components'

const props = defineProps<{
  name: string,
  service,
  form,
}>()
const emit = defineEmits<{
  (e: "refresh"): void,
}>()

const modal = ref(null)

function show() {
  console.log(`show()!`)
  modal.value.show()
}

</script>

<template>
  <n-button type="success" @click="show">
    <template #icon><n-icon :component="Add" /></template>
    Add
  </n-button>
  <s-modal-create
    ref="modal"
    :name="name"
    :service="service"
    :form="form"
    @refresh="emit('refresh')"
  />
</template>
