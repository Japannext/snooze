<script setup lang="ts">

import { computed, defineEmits } from 'vue'
import { NSelect, NInputGroup, NAutoComplete } from 'naive-ui'
import type { SelectOption } from 'naive-ui'

interface Props {
  value: string[],
  size: string,
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: "update:value", value: string[]): void,
}>()

const dataValue = computed({
  get() { return props.value
  },
  set(v) {
    emit("update:value", v)
  },
})

const nestedKeys: string[] = [
  "resource",
  "attributes",
  "severity",
]

const coreOptions: SelectOption[] = [
  {value: "resource", label: "resource"},
  {value: "attributes", label: "attributes"},
  {value: "severity", label: "severity"},
]

const resourceOptions: string[] = [
  "host.hostname",
  "k8s.cluster.name",
  "k8s.container.name",
  "k8s.cronjob.name",
  "k8s.daemonset.name",
  "k8s.deployment.name",
  "k8s.job.name",
  "k8s.namespace.name",
  "k8s.node.name",
  "k8s.pod.name",
  "k8s.replicaset.name",
  "k8s.statefulset.name",
  "service.name",
  "service.version",
  "event.name",
  "event.domain",
]

const attributeOptions: string[] = [
  "exception.message",
  "exception.stackstrace",
  "exception.type",
  "feature_flag.key",
  "feature_flag.provider_name",
  "feature_flag.variant",
  "log.file.name",
  "log.file.name_resolved",
  "log.file.path",
  "log.file.path_resolved",
  "log.iostream",
  "net.app.protocol.name",
  "net.app.protocol.version",
  "net.host.carrier.icc",
  "net.host.carrier.mcc",
  "net.host.carrier.mnc",
  "net.host.carrier.name",
  "net.host.connection.subtype",
  "net.host.connection.type",
  "net.host.name",
  "net.host.port",
  "net.peer.name",
  "net.peer.port",
  "net.sock.family",
  "net.sock.host.addr",
  "net.sock.host.port",
  "net.sock.peer.addr",
  "net.sock.peer.name",
  "net.sock.peer.port",
  "net.transport",
  "syslog.facility",
  "syslog.format",
  "syslog.msgid",
  "syslog.procid",
]

const severityOptions: string[] = [
  "text",
  "number",
]

const options = computed(() => {
  switch(dataValue.value[0]) {
    case "resource":
      return resourceOptions
    case "attributes":
      return attributeOptions
    case "severity":
      return severityOptions
    default:
      return []
  }
})
</script>

<template>
  <n-select
    v-model:value="dataValue[0]"
    style="width: 100px;"
    :consistent-menu-width="false"
    :size="size"
    :options="coreOptions"
  />
  <n-auto-complete
    v-if="dataValue.length > 0 && nestedKeys.includes(dataValue[0])"
    v-model:value="dataValue[1]"
    placeholder="Key"
    style="min-width: 200px;"
    :consistent-menu-width="false"
    :size="size"
    :options="options"
  />
</template>
