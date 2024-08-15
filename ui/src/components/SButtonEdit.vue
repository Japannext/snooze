<script setup lang="ts">
import { ref, defineEmits, withDefaults } from 'vue'
import { NButton, NIcon } from 'naive-ui'
import { Pencil } from '@vicons/ionicons5'

import { SModalEdit } from '@/components'

const props = withDefaults(defineProps<{
  name: string,
  service,
  form,
  withModal: boolean,
}>(), {
  withModal: false,
})

const emit = defineEmits<{
  (e: "refresh"): void,
}>()

const modal = ref(null)

function show() {
  modal.value.show()
}

</script>

<template>
  <n-button type="warning" round @click="show" size="small">
    <template #icon><n-icon :component="Pencil" /></template>
  </n-button>
  <s-modal-edit
    v-if="withModal"
    ref="modal"
    :name="name"
    :service="service"
    :form="form"
    @refresh="emit('refresh')"
  />
</template>
