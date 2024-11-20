<script setup lang="ts">
import { defineProps } from 'vue'

import { NSpace, NTag, NIcon } from 'naive-ui'
import { XTag } from '@/components/attributes'
import { TagDismiss16Regular } from '@/icons'

import type { Tag } from '@/api'

defineProps<{
  tags?: Tag[],
}>()

</script>

<template>
  <n-space v-if="tags" size="small">
    <x-tag
      v-for="tag in tags"
      :key="tag.name"
      :tag="tag"
    />
  </n-space>
  <n-tag v-else checkable disabled>
    <template #icon><n-icon :component="TagDismiss16Regular" /></template>
  </n-tag>
</template>
