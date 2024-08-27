<script setup lang="ts">
import { defineProps, computed } from 'vue'
import { NTag, NIcon } from 'naive-ui'
import { Alert, Mail, LogoGoogle } from '@vicons/ionicons5'
import { BuildingLighthouse20Filled } from '@vicons/fluent'

import type { Destination } from '@/api/types'

const props = defineProps<{
  destination: Destination,
}>()

const component = computed(() => {
  switch(props.destination.queue) {
    case "mail":
      return Mail
    case "googlechat":
      return LogoGoogle
    case "patlite":
      return BuildingLighthouse20Filled
    default:
      return Alert
  }
})

</script>

<template>
  <div>
    <n-tag>
      <template #icon><n-icon :component="component" /></template>
      {{ destination.profile }}
    </n-tag>
  </div>
</template>
