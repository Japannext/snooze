<script setup lang="ts">
import { defineProps } from 'vue'
import { NDescriptions, NDescriptionsItem, NTag, NIcon, NSpace, NPopover } from 'naive-ui'

import { Alert } from '@vicons/ionicons5'
import { Server, Terminal, Hashtag } from '@vicons/fa'
import { ContainerSoftware, Kubernetes } from '@vicons/carbon'
import { ArrowForwardIosFilled } from '@vicons/material'

defineProps<{
  identity: object|undefined,
}>()

function shortname(host: string): string {
  if (host.includes(".")) {
    return host.split('.')[0]
  }
  return host
}

function k8sKind(kind: string): string {
  var subKind = kind.replace("k8s.", "")
  switch (subKind) {
    case "deployment":
      return "deploy"
    case "replicaset":
      return "rs"
    case "statefulset":
      return "sts"
    case "persistentvolumeclaim":
      return "pvc"
  }
  return subKind
}

</script>

<template>
  <template v-if="identity">
    <template v-if="identity.kind == 'host'">
      <n-tag v-if="identity.host" size="small">
        <template #icon><n-icon :component="Server" /></template>
        {{ shortname(identity.host) }}<b v-if="identity.process"># {{ identity.process }}</b>
      </n-tag>
    </template>
    <n-popover v-else-if="identity.kind && (identity.kind).startsWith('k8s.')" trigger="click">
      <template #trigger>
        <n-tag size="small">
          <template #icon>
            <n-icon :component="Kubernetes" />
          </template>
          {{ k8sKind(identity.kind) }}/{{ identity.name }} <b v-if="identity.namespace">@{{ identity.namespace }}</b>
        </n-tag>
      </template>
      <template #default>
        <n-icon :component="Kubernetes" size="40" />
        <n-descriptions size="small" label-placement="left">
          <n-descriptions-item label="Kind">
            {{ identity.kind }}
          </n-descriptions-item>
          <n-descriptions-item label="Name">
            {{ identity.name }}
          </n-descriptions-item>
          <n-descriptions-item label="Namespace">
            {{ identity.namespace }}
          </n-descriptions-item>
        </n-descriptions>
      </template>
    </n-popover>
  </template>
</template>
