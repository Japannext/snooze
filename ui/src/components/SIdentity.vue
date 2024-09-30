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

</script>

<template>
  <template v-if="!props.identity || !props.identity.kind">
    {{ identity }}
  </template>
  <template v-else-if="props.identity.kind == 'host'">
    <n-tag size="small">
      <template #icon><n-icon :component="Server" /></template>
      {{ identity.host.split('.')[0] }}
    </n-tag>
    <n-tag v-if="props.identity.process" size="small">
      <template #icon><n-icon :component="Hashtag" /></template>
      {{ identity.process }}
    </n-tag>
  </template>
  <template v-else-if="(props.identity.kind).startsWith('k8s.')">
    <n-tag size="small">
      <template #icon><n-icon :component="ContainerSoftware" /></template>
      {{ identity.name }} <b v-if="identity.namespace">@{{ identity.namespace }}</b>
    </n-tag>
  </template>
</template>
