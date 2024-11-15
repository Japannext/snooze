<script setup lang="ts">
import { defineProps, computed } from 'vue'

import { NTag } from 'naive-ui'

const props = defineProps<{
  kind: string,
}>()

const KINDS = {
  "active": {type: "", value: "active"},
  "snoozed": {type: "warning", value: "snoozed"},
  "silenced": {type: "info", value: "silenced"},
  "acked": {type: "success", value: "acked"},
}

const tag = computed(() => {
  var k = KINDS[props.kind]
  if (k) {
    return k
  }
  return {type: "", value: "-"}
})

</script>

<template>
  <n-tag :type="tag.type">
    {{ tag.value }}
  </n-tag>
</template>
