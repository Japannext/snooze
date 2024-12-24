<script setup lang="ts">
import { defineProps } from 'vue'
import { NSpace } from 'naive-ui'

import type { Alert } from '@/api'

import { XAlertName, XIdentity, XSeverity } from '@/components/attributes'

// props
defineProps<{
  row: Alert,
}>()

</script>

<template>
  <n-space size="small">
    <!-- <s-alert-group v-if="row.alertGroup" :value="row.alertGroup" /> -->
    <x-identity v-if="row.identity" :identity="row.identity" />
    <x-alert-name v-if="row.alertName" :value="row.alertName" />
    <x-severity v-if="row.severityNumber !== undefined" :text="row.severityText" :number="row.severityNumber" />
  </n-space>
</template>
