<script setup lang="ts">
import { defineProps, computed } from 'vue'

import { NTag } from 'naive-ui'

import type { LogStatus, LogStatusKind } from '@/api'

const props = defineProps<{
  status: LogStatus,
}>()

const KINDS = {
  0: {type: "", value: "active"},
  1: {type: "warning", value: "snoozed"},
  2: {type: "info", value: "silenced"},
  3: {type: "error", value: "ratelimited"}, // SHOULD NOT HAPPEN
  4: {type: "error", value: "dropped"},     // SHOULD NOT HAPPEN
  5: {type: "error", value: "activecheck"}, // SHOULD NOT HAPPEN
  6: {type: "success", value: "acked"},
}

const tag = computed(() => {
  if (props.status) {
    var k = KINDS[props.status.kind]
    if (k) {
      return k
    }
    return {type: "", value: "???"}
  }
  return {type: "", value: "???"}
})

</script>

<template>
  <n-tag :type="tag.type">
    {{ tag.value }}
  </n-tag>
</template>
