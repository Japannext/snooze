<script setup lang="ts">
import { ref, defineProps, onMounted } from 'vue'
import { NSpace } from 'naive-ui'

import type { Log } from '@/api/types'

import SHostname from '@/components/attributes/SHostname.vue'
import SProcess from '@/components/attributes/SProcess.vue'
import SLogPattern from '@/components/attributes/SLogPattern.vue'
import SSeverity from '@/components/attributes/SSeverity.vue'

const props = defineProps<{
  row: Log,
}>()
</script>

<template>
  <n-space size="small">
    <s-hostname v-if="row.identity.kind == 'host'" :hostname="row.identity.hostname" />
    <template v-if="row.pattern">
      <s-log-pattern :name="row.pattern" />
    </template>
    <template v-else>
      <s-process v-if="row.identity.process" :process="row.identity.process" />
    </template>
    <s-severity :text="row.severityText" :number="row.severityNumber" />
  </n-space>
</template>
