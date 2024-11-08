<script setup lang="ts">
import { defineProps, defineModel } from 'vue'

// Components
import { NTab, NTabs } from 'naive-ui'

defineProps<{
  filters: Array<object>,
}>()

const model = defineModel()

</script>

<template>
  <n-tabs v-model:value="model" type="line">
    <n-tab v-for="filter in filters" :key="filter.value" :name="filter.value">
      {{ filter.label }}
    </n-tab>
  </n-tabs>
</template>
