<script setup lang="ts">
const props = defineProps<{
  fieldName: string;
  title: string;
}>();
</script>

<template>
  <Column :field="fieldName" :header="title" style="width: 50px">
    <template #body="{ data }">
      <Tag
        :value="data[fieldName]"
        :style="{ 'background-color': data[fieldName], color: 'white' }"
      />
    </template>
  </Column>
</template>
