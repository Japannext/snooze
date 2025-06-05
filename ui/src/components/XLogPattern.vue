<script setup lang="ts">
import { ref } from "vue";
import { TypePattern, BookOpen } from "@/icons";

const props = defineProps<{
  name: string;
}>();

const popover = ref();

function toggle(event) {
  if (popover.value) {
    popover.value.toggle(event);
  }
}
</script>

<template>
  <Button
    size="small"
    severity="secondary"
    icon="pi pi-objects-column"
    :label="props.name"
    @click="toggle"
    rounded
  />
  <Popover ref="popover"> test </Popover>
</template>
