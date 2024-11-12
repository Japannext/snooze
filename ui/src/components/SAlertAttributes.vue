<script setup lang="ts">
import { defineProps } from 'vue'
import { NSpace } from 'naive-ui'

import type { Log } from '@/api/types'

import SAlertGroup from '@/components/attributes/SAlertGroup.vue'
import SAlertName from '@/components/attributes/SAlertName.vue'
import SHostname from '@/components/attributes/SHostname.vue'
import SProcess from '@/components/attributes/SProcess.vue'
import SLogPattern from '@/components/attributes/SLogPattern.vue'
import SIdentity from '@/components/attributes/SIdentity.vue'

import { XSeverity } from '.'

// props
defineProps<{
  row: Log,
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
