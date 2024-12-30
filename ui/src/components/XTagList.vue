<script setup lang="ts">
import { defineProps, computed } from 'vue'

import { NSpace, NTag, NIcon } from 'naive-ui'
import { XTag } from '@/components'
import { TagDismiss16Regular } from '@/icons'
import type { Tag } from '@/api'
import { tagMap } from '@/stores'

const props = defineProps<{
  tags: Tag[],
}>()
</script>

<template>
  <n-space v-if="tags" size="small">
    <x-tag
      v-for="tag in tags"
      v-bind:key="tag.name"
      :name="tag.name"
      :color="tag.color"
    />
  </n-space>
  <n-tag v-else checkable disabled>
    <template #icon><n-icon :component="TagDismiss16Regular" /></template>
  </n-tag>
</template>
