<script setup lang="ts">
import { NSpace, NEllipsis } from 'naive-ui'
import { XTag } from '@/components'

defineProps<{
  name: string,
  color?: string,
  description?: string,
}>()
</script>

<template>
  <n-space inline>
    <x-tag :name="name" :color="color" />
    <n-ellipsis v-if="description" style="max-width: 500px" :tooltip="false">
      {{ description }}
    </n-ellipsis>
  </n-space>
</template>
