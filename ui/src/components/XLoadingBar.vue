<script setup lang="ts">
interface Props {
  loading: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
});
</script>

<template>
  <div style="height: 2px">
    <ProgressBar v-if="loading" mode="indeterminate" style="height: 2px" />
  </div>
</template>
