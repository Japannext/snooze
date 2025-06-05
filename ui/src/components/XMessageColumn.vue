<script setup lang="ts">
import { type Log } from "@/api";
</script>

<template>
  <Column field="message" header="Message">
    <template #body="{ data }">
      <x-log-pattern v-if="data.pattern" :name="data.pattern" />
      <Tag
        v-else-if="'process' in data.identity"
        severity="secondary"
        icon="pi pi-hashtag"
        rounded
      >
        {{ data.identity.process }}
      </Tag>
      <span>{{ data.message }}</span>
    </template>
  </Column>
</template>
