<script setup lang="ts">
import { NSpace, NTag, NIcon } from 'naive-ui'
import { TagDismiss16Regular } from '@/icons'
import type { Group } from '@/api'

defineProps<{
  groups: Group[],
}>()

function printLabels(gr: Group): string {
  if (gr.labels) {
    return Object.keys(gr.labels).map(key => gr.labels[key]).join(',')
  }
  return ""
}
</script>

<template>
  <n-space v-if="groups" size="small">
    <n-tag
      v-for="group in groups"
      :key="group.hash"
    >
      [{{ group.name }}] {{ printLabels(group) }}
    </n-tag>
  </n-space>
  <n-tag v-else checkable disabled>
    <template #icon><n-icon :component="TagDismiss16Regular" /></template>
  </n-tag>
</template>
