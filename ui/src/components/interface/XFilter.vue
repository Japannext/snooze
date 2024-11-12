<script setup lang="ts">
import { defineProps, defineModel } from 'vue'

// Components
// import { NTab, NTabs } from 'naive-ui'
import { NInputGroup, NInputGroupLabel, NRadioGroup, NRadioButton } from 'naive-ui'

defineProps<{
  filters: Array<object>,
}>()

const model = defineModel('value')

</script>

<template>
  <!--
  <n-tabs v-model:value="model" type="line">
    <n-tab v-for="filter in filters" :key="filter.value" :name="filter.value">
      {{ filter.label }}
    </n-tab>
  </n-tabs>
  -->
  <n-input-group>
    <!--<n-input-group-label>Filters</n-input-group-label>-->
    <n-radio-group v-model:value="model">
      <n-radio-button
        v-for="item in filters"
        :key="item.value"
        :label="item.label"
        :value="item.value"
      />
    </n-radio-group>
  </n-input-group>
</template>
