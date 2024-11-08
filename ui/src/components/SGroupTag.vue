<script setup lang="ts">
import { defineProps } from 'vue'
import type { Group } from "@/api/types"

defineProps<{
  group: Group,
}>()

</script>

<template>
  <span>
    {{ group.name }} {{ group.labels }}<br />
  </span>
</template>
