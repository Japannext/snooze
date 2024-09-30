<script setup lang="ts">
import { defineProps } from 'vue'
import { NTag } from 'naive-ui'

const props = defineProps<{
  value: string,
}>()

</script>
<template>
  <n-tag size="small">
    {{ props.value }}
  </n-tag>
</template>
