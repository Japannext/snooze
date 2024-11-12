<script setup lang="ts">
import { ref, defineProps, onMounted } from 'vue'
import { NSpace } from 'naive-ui'

import type { Log } from '@/api/types'
import { XLogPattern, XIdentity, XSeverity } from '@/components/attributes'

defineProps<{
  row: Log,
}>()
</script>

<template>
  <n-space size="small">
    <x-identity :identity="row.identity" />
    <template v-if="row.pattern">
      <x-log-pattern :name="row.pattern" />
    </template>
    <x-severity :text="row.severityText" :number="row.severityNumber" />
  </n-space>
</template>
