<script setup lang="ts">
import { defineProps } from "vue";
import { NTag } from "naive-ui";

const props = defineProps<{
  value: string;
}>();
</script>
<template>
  <Button :label="props.value" size="small" severity="secondary" />
</template>
