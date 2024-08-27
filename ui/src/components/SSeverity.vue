<script setup lang="ts">
import { computed } from 'vue'
import { NTag } from 'naive-ui'

const props = defineProps<{
  text: string,
  number: number,
  size: string,
}>()
</script>

<template>
  <n-tag :size="size">{{ props.text || "---" }}</n-tag>
</template>
