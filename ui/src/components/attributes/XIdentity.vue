<script setup lang="ts">
import { defineProps } from 'vue'
import { NList, NListItem, NTag, NIcon, NSpace, NPopover } from 'naive-ui'

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
    <n-popover v-else-if="identity.kind && (identity.kind).startsWith('k8s.')" trigger="click" style="padding: 0;">
      <template #trigger>
        <n-tag size="small">
          <template #icon>
            <n-icon :component="Kubernetes" />
          </template>
          {{ k8sKind(identity.kind) }}/{{ identity.name }} <b v-if="identity.namespace">@{{ identity.namespace }}</b>
        </n-tag>
      </template>
      <template #default>
        <n-list bordered>
          <template #header>
            <n-space justify="center">
            <n-icon :component="Kubernetes" />
            Kubernetes resource
            </n-space>
          </template>
          <n-list-item>
            <b>Kind</b>: {{ identity.kind }}
          </n-list-item>
          <n-list-item>
            <b>Name</b>: {{ identity.name }}
          </n-list-item>
          <n-list-item>
            <b>Namespace</b>: {{ identity.namespace }}
          </n-list-item>
        </n-list>
      </template>
    </n-popover>
  </template>
</template>

<style>
.n-tag {
  cursor: help;
}
</style>
