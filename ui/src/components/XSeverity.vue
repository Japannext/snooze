<script setup lang="ts">
import { defineProps, ref, onMounted } from "vue";

const props = defineProps<{
  text: string | undefined;
  number: number;
}>();

const color = ref<string>("default");

onMounted(() => {
  if (props.number == 0) {
    // Not set
    color.value = "default";
  } else if (props.number < 5) {
    // TRACE
    color.value = "default";
  } else if (props.number < 9) {
    // DEBUG
    color.value = "default";
  } else if (props.number < 13) {
    // INFO
    color.value = "info";
  } else if (props.number < 17) {
    // WARN
    color.value = "warning";
  } else if (props.number < 21) {
    // ERROR
    color.value = "danger";
  } else if (props.number < 25) {
    // FATAL
    color.value = "danger";
  } else {
    // UNKNOWN
    color.value = "default";
  }
});
</script>
<template>
  <Tag v-if="text" :severity="color" :value="text" />
</template>
