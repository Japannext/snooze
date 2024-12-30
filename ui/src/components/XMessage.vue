<script setup lang="ts">
import { type Log } from '@/api'

import { NTag, NIcon, NSpace, NPerformantEllipsis } from 'naive-ui'
import { XLogPattern } from '@/components'
import { Hashtag } from '@/icons'

defineProps<{
  log: Log,
}>()

</script>

<template>
  <n-space inline align="center">
    <x-log-pattern v-if="log.pattern" :name="log.pattern" />
    <n-tag v-else-if="'process' in log.identity" size="small">
      <template #icon><n-icon :component="Hashtag" /></template>
      {{ log.identity.process }}
    </n-tag>
    <span>{{ log.message }}</span>
  </n-space>
</template>
