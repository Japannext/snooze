<script setup lang="ts">

import { computed } from 'vue'

import {
  NButton,
} from 'naive-ui'

import { Severity } from '@/api'

interface Props {
  severity: Severity,
  size: string,
}

const props = defineProps<Props>()

const buttonAttributes = computed(() => {
  switch(props.severity.number) {
    // TRACE
    case 1: case 2: case 3: case 4:
      return {type: "success", dashed: true}
    // DEBUG
    case 5: case 6: case 7: case 8:
      return {type: "info", dashed: true}
    // INFO
    case 9: case 10: case 11: case 12:
      return {type: "info", secondary: true}
    // WARNING
    case 13: case 14: case 15: case 16:
      return {type: "warning", secondary: true}
    // ERROR
    case 17: case 18: case 19: case 20:
      return {type: "error", secondary: true}
    // FATAL
    case 21: case 22: case 23: case 24:
      return {type: "error"}
    // UNKNOWN
    default:
      return {type: "warning", dashed: true}
  }
})
</script>

<template>
  <n-button v-bind="buttonAttributes" :size="size">{{ severity.text || "---" }}</n-button>
</template>
