<script setup lang="ts">
import { defineProps } from "vue";
import type { Alert } from "@/api";
</script>

<template>
  <Column header="Attributes">
    <template #body="{ data }">
      <XIdentity v-if="data.identity" :identity="data.identity" />
      <XAlertName v-if="data.alertName" :value="data.alertName" />
      <XSeverity
        v-if="data.severityNumber !== undefined"
        :text="data.severityText"
        :number="data.severityNumber"
      />
    </template>
  </Column>
</template>
