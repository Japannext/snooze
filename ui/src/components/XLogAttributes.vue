<script setup lang="ts">
import { defineProps } from 'vue'
import { NSpace } from 'naive-ui'

import type { Log } from '@/api'
import { XLogPattern, XIdentity, XSeverity } from '@/components'

defineProps<{
  row: Log,
}>()
</script>

<template>
  <n-space size="small">
    <x-identity :identity="row.identity" />
    <x-severity :text="row.severityText" :number="row.severityNumber" />
  </n-space>
</template>
