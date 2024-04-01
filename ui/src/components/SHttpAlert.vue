<script setup lang="ts">
// Display an HTTP alert, with rendering adapted to response code.

import { ref, computed, defineEmits, withDefaults } from 'vue'
import { NAlert } from 'naive-ui'
import type { Ref } from 'vue'

import { HttpAlert } from '@/utils/validation'

interface Props {
  // The HTTP alert
  alert: Ref<HttpAlert>, // TODO: Remove the reference
}

const props = withDefaults(defineProps<Props>(), {
  alert: () => { return {} },
})

const emit = defineEmits<{
  (e: "update:closed"): void,
}>()

const variant = computed(() => {
  if (200 <= props.alert.value.code && props.alert.value.code < 400) {
    return "success"
  } else if (400 <= props.alert.value.code && props.alert.value.code <= 599) {
    return "error"
  } else {
    return "info"
  }
})

function close() {
  emit("update:closed")
}

</script>

<template>
  <n-alert
    v-if="alert.value.show"
    :title="`HTTP ${alert.value.code}: ${alert.value.name}`"
    :type="variant"
    closable
    :on-after-leave="close"
  >
    {{ alert.value.text }}
  </n-alert>
</template>
