<script setup lang="ts">
import { NSpace, NTag } from 'naive-ui'

const props = defineProps<{
  value: Array<string>,
}>()


</script>

<template>
  <n-space>
    <n-tag v-for="text in value" :key="text">
      {{ text }}
    </n-tag>
  </n-space>
</template>
