<script setup lang="ts">
import { defineProps } from 'vue'
import { NTag, NIcon } from 'naive-ui'
import { Hashtag } from '@/icons'

const props = defineProps<{
  process: string,
}>()

</script>
<template>
  <n-tag size="small">
    <template #icon>
      <n-icon :component="Hashtag" />
    </template>
    {{ props.process }}
  </n-tag>
</template>
