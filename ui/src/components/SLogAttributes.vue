<script setup lang="ts">
import { ref, defineProps, onMounted } from 'vue'
import { NSpace } from 'naive-ui'

import type { Log } from '@/api/types'
import SLogPattern from '@/components/attributes/SLogPattern.vue'
import SSeverity from '@/components/attributes/SSeverity.vue'
import SIdentity from '@/components/attributes/SIdentity.vue'

const props = defineProps<{
  row: Log,
}>()
</script>

<template>
  <n-space size="small">
    <s-identity :identity="row.identity" />
    <template v-if="row.pattern">
      <s-log-pattern :name="row.pattern" />
    </template>
    <s-severity :text="row.severityText" :number="row.severityNumber" />
  </n-space>
</template>
