<script setup lang="ts">
import { defineProps, defineModel } from "vue";
import type { Filter } from "@/api";

defineProps<{
  filters: Array<Filter>;
}>();

const model = defineModel<string>("value");
</script>

<template>
  <!--
  <n-tabs v-model:value="model" type="line">
    <n-tab v-for="filter in filters" :key="filter.value" :name="filter.value">
      {{ filter.label }}
    </n-tab>
  </n-tabs>
  -->
  <SelectButton
    v-model="model"
    :options="filters"
    size="small"
    :allowEmpty="false"
    optionLabel="label"
    optionValue="value"
  />
  <!-- <n-input-group> -->
  <!--<n-input-group-label>Filters</n-input-group-label>-->
  <!-- <n-radio-group v-model:value="model"> -->
  <!-- <n-radio-button -->
  <!-- v-for="item in filters" -->
  <!-- :key="item.value" -->
  <!-- :label="item.label" -->
  <!-- :value="item.value" -->
  <!-- /> -->
  <!-- </n-radio-group> -->
  <!-- </n-input-group> -->
</template>
