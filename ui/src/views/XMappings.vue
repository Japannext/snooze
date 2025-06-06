<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getPipelineMappings, type Mapping } from "@/api";

const loading = ref<boolean>(false);
const selected = ref<Mapping>({});
const options = ref<string[]>([]);

function refresh() {
  loading.value = true;
  getPipelineMappings()
    .then((list) => {
      options.value = list;
    })
    .finally(() => {
      loading.value = false;
    });
}

onMounted(() => {
  refresh();
});

function getMappings() {
  return Object.entries(selected.value.map);
}
</script>

<template>
  <div style="height: 2px">
    <ProgressBar v-if="loading" mode="indeterminate" style="height: 2px" />
  </div>
  <Select v-model="selected" :options="options" optionLabel="name" />
  {{ selected }}
  <DataTable v-if="selected" :value="getMappings()">
    <Column field="key" header="Key" />
    <Column field="value" header="Value" />
  </DataTable>
</template>
