<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getPipelineGroupings, type Grouping } from "@/api";

const loading = ref<boolean>(false);
const items = ref<Grouping[]>([]);

function refresh() {
  loading.value = true;
  getPipelineGroupings()
    .then((list) => {
      items.value = list;
    })
    .finally(() => {
      loading.value = false;
    });
}

onMounted(() => {
  refresh();
});
</script>

<template>
  <div style="height: 2px">
    <ProgressBar v-if="loading" mode="indeterminate" style="height: 2px" />
  </div>
  <DataTable :value="items" size="small" stripedRows showGridlines>
    <Column field="name" header="Name" />
    <Column field="if" header="If" />
  </DataTable>
</template>
