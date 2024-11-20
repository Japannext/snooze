<script setup lang="ts">
import { defineProps, onMounted } from 'vue'

import { NSpace, NEllipsis } from 'naive-ui'
import { XTag } from '@/components/attributes'

import type { Tag } from '@/api'

const props = defineProps<{
  tag: Tag,
}>()

</script>

<template>
  <n-space inline>
    <x-tag :tag="tag" />
    <n-ellipsis v-if="tag.description" style="max-width: 500px" :tooltip="false">
      {{ tag.description }}
    </n-ellipsis>
  </n-space>
</template>
