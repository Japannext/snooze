<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps<{
  identity: object,
}>()

</script>

<template>
  <span>
    <template v-if="identity.kind == 'host'">
      {{ identity.hostname }} # {{ props.identity.process }}
    </template>
    <template v-else-if="identity.kind == 'k8s.pod'">
      {{ identity.pod }} @{{ props.identity.namespace }}
    </template>
    <template v-else>
      {{ identity }}
    </template>
  </span>
</template>
