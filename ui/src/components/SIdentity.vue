<script setup lang="ts">
import { defineProps, computed } from 'vue'
import { NTag, NIcon, NSpace } from 'naive-ui'

import { Alert } from '@vicons/ionicons5'
import { Server, Terminal, Hashtag } from '@vicons/fa'
import { ContainerSoftware } from '@vicons/carbon'
import { ArrowForwardIosFilled } from '@vicons/material'

const props = defineProps<{
  identity: object,
}>()

const component = computed(() => {
  switch(props.identity.kind) {
    case "host":
      return Server
    case "pod":
      return ContainerSoftware
    default:
      return Alert
  }
})

const text = computed(() => {
  switch(props.identity.kind) {
    case "host":
      return `${props.identity.hostname}#${props.identity.process}`
    case "pod":
      return `${props.identity.pod}@${props.identity.namespace}`
    default:
      return JSON.stringify(props.identity)
  }
})
</script>

<template>
  <template v-if="props.identity.kind == 'host'">
    <n-tag size="small">
      <template #icon><n-icon :component="Server" /></template>
      {{ identity.hostname.split('.')[0] }}
    </n-tag>
    <n-tag v-if="props.identity.process" size="small">
      <template #icon><n-icon :component="Hashtag" /></template>
      {{ identity.process }}
    </n-tag>
  </template>
</template>
