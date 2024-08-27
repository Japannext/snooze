<script setup lang="ts">
import { defineProps } from 'vue'
import { NTag, NIcon } from 'naive-ui'
import { Server } from '@/icons'

const props = defineProps<{
  hostname: string,
}>()

</script>
<template>
  <n-tag size="small">
    <template #icon>
      <n-icon :component="Server" />
    </template>
    {{ props.hostname.split('.')[0] }}
  </n-tag>
</template>
