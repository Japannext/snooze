<script setup lang="ts">
interface Props {
  fieldName: string;
  label: string;
}

const props = defineProps<Props>();
</script>

<template>
  <Column :field="props.fieldName" :header="props.label">
    <template #body="{ data }">
      <pre>{{ data[fieldName] }}</pre>
    </template>
  </Column>
</template>
