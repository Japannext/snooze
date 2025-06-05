<script setup lang="ts">
import { defineProps } from "vue";
import type { Identity, HostIdentity } from "@/api";

import { Server, Hashtag } from "@vicons/fa";
import { Kubernetes } from "@vicons/carbon";

defineProps<{
  identity: Identity;
}>();

function shortname(host: string): string {
  if (host.includes(".")) {
    return host.split(".")[0];
  }
  return host;
}

function k8sShortname(component: string): string {
  switch (component) {
    case "deployment":
      return "deploy";
    case "statefulset":
      return "sts";
    case "daemonset":
      return "ds";
    case "replicaset":
      return "rs";
    default:
      return component;
  }
}
</script>

<template>
  <span>
    <Tag v-if="identity.host" icon="pi pi-server" severity="secondary">
      {{ shortname(identity.host) }}
    </Tag>

    <Tag v-else-if="identity.kind == 'kubernetes'" icon="pi pi-android">
      {{ k8sShortname(identity.component) }}/{{ identity.name }}
      <b>@{{ identity.namespace }}</b>
    </Tag>
  </span>
</template>
