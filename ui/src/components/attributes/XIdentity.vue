<script setup lang="ts">
import { defineProps } from 'vue'
import { NList, NListItem, NTag, NIcon, NSpace, NPopover } from 'naive-ui'
import type { Identity, HostIdentity } from '@/api'

import { Server, Hashtag } from '@vicons/fa'
import { Kubernetes } from '@vicons/carbon'

defineProps<{
  identity: Identity,
}>()

function shortname(host: string): string {
  if (host.includes(".")) {
    return host.split('.')[0]
  }
  return host
}

function k8sShortname(component: string): string {
  switch(component) {
    case "deployment":
      return "deploy"
    case "statefulset":
      return "sts"
    case "daemonset":
      return "ds"
    case "replicaset":
      return "rs"
    default:
      return component
  }

}

</script>

<template>
  <template v-if="identity.kind == 'host'">
    <n-tag v-if="identity.hostname" size="small">
      <template #icon><n-icon :component="Server" /></template>
      {{ shortname(identity.hostname) }}<b v-if="identity.process"># {{ identity.process }}</b>
    </n-tag>
  </template>
  <template v-else-if="identity.kind == 'kubernetes'">
    <n-tag size="small">
      <template #icon><n-icon :component="Kubernetes" /></template>
      {{ k8sShortname(identity.component) }}/{{ identity.name }} <b>@{{ identity.namespace }}</b>
    </n-tag>
  </template>
</template>
