<script setup lang="ts">
import { defineProps } from 'vue'
import { NSpace } from 'naive-ui'

import type { Alert } from '@/api'

import { XAlertGroup, XAlertName, XHostname, XProcess, XIdentity, XSeverity } from '@/components/attributes'

// props
defineProps<{
  row: Alert,
}>()

</script>

<template>
  <n-space size="small">
    <!-- <s-alert-group v-if="row.alertGroup" :value="row.alertGroup" /> -->
    <s-identity v-if="row.identity" :identity="row.identity" />
    <s-alert-name v-if="row.alertName" :value="row.alertName" />
    <x-severity :text="row.severityText" :number="row.severityNumber" />
  </n-space>
</template>
