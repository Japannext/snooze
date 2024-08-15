<script setup lang="ts">

import { ref, defineEmits, onMounted, VNodeChild } from 'vue'

import {
  NSpace,
  NIcon, NTag,
  NButton,
  NPopover, NTooltip,
} from 'naive-ui'

import { Server, Boat } from '@vicons/ionicons5'

interface Props {
  resource: object,
  size: string,
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: "select", filter: Map<string, string>): void,
}>()

// Kind
enum Kind {
  K8s = "k8s",
  Syslog = "syslog",
  Unknown = "unknown",
}

function detectKind(res: object): Kind {
  if (Object.hasOwn(res, "k8s.cluster.name")) {
    return Kind.K8s
  } else if (Object.hasOwn(res, "host.hostname")) {
    return Kind.Syslog
  }
    return Kind.Unknown
}

// Kubernetes Objects
interface K8sObject {
  name: string
  short?: string
}

interface K8sInstance extends K8sObject {
  instance: string
}

const k8sObjects = new Map<string, K8sObject>([
  ["k8s.deployment.name", {name: "deployment", short: "deploy"}],
  ["k8s.statefulset.name", {name: "statefulset", short: "sts"}],
  ["k8s.daemonset.name", {name: "daemonset", short: "ds"}],
  ["k8s.job.name", {name: "job", short: "job"}],
  ["k8s.cronjob.name", {name: "cronjob", short: "cj"}],
  ["k8s.replicaset.name", {name: "replicaset", short: "rs"}],
  ["k8s.pod.name", {name: "pod", short: "pod"}],
])

function detectK8sInstance(res: object): K8sInstance {
  for (let [key, value] of k8sObjects) {
    if (Object.hasOwn(res, key)) {
      return {instance: res[key], ...value}
    }
  }
  return {name: "unknown", instance: "unknown"}
}

const kind = ref<Kind>(null)
const k8sObj = ref<K8sInstance>(null)
const resource = ref<object>({})

// Render
onMounted(() => {
  resource.value = props.resource
  kind.value = detectKind(props.resource)
  if (kind.value == Kind.K8s) {
    k8sObj.value = detectK8sInstance(props.resource)
  }
})

</script>

<template>
  <n-space :size="2">
    <template v-if="kind == 'k8s'">
      <n-button secondary type="info" :size="size">
        <template #icon><n-icon :component="Boat" /></template>
        {{ k8sObj.short }} / {{ k8sObj.instance }}
      </n-button>
      <n-button
        v-if="Object.hasOwn(resource, 'k8s.container.name')"
        type="info"
        :size="size"
        secondary
        round
      >
        {{ resource["k8s.container.name"] }}
      </n-button>
    </template>

    <template v-else-if="kind == 'syslog'">
      <n-button type="success" secondary :size="size">
        <template #icon><n-icon :component="Server" /></template>
        {{ resource["host.hostname"] }}
      </n-button>
      <n-button
        v-if="Object.hasOwn(resource, 'service.name')"
        type="success"
        :size="size"
        secondary
        round
      >
        {{ resource["service.name"] }}
      </n-button>
    </template>
  </n-space>
</template>
