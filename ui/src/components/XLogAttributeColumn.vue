<script setup lang="ts">
import { defineProps } from "vue";
import type { Alert } from "@/api";
</script>

<template>
  <Column header="Attributes" style="width: 300px">
    <template #body="{ data }">
      <XIdentity v-if="data.identity" :identity="data.identity" />
      <XSeverity
        v-if="data.severityNumber !== undefined"
        :text="data.severityText"
        :number="data.severityNumber"
      />
    </template>
  </Column>
</template>

<style lang="scss" scoped>
.p-tag {
  margin: 0 5px;
}
</style>
