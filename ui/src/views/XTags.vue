<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouteQuery } from "@vueuse/router";
import { useToast } from "primevue/usetoast";

import { getTags, type Tag, type GetTagsParams } from "@/api";

const items = ref<Tag[]>();
const selectedItems = ref<string[]>([]);

const loading = ref<boolean>(false);
const toast = useToast();
const showNewTagModal = ref<boolean>(false);
const showDeleteTagModal = ref<boolean>(false);

const page = useRouteQuery<number>("page", 1, { transform: Number });
const size = useRouteQuery<number>("size", 20, { transform: Number });
const total = ref<number>(0);
const more = ref<boolean>(false);

const params = ref<GetTagsParams>({
  search: "",
  pagination: {},
});

onMounted(() => {
  refresh();
});

function refresh() {
  loading.value = true;
  params.value.pagination = {
    page: page.value,
    pageSize: size.value,
  };
  getTags(params.value)
    .then((list) => {
      items.value = list.items;
      total.value = list.total;
      more.value = list.more;
    })
    .catch((err) => {
      items.value = [];
      toast.add({
        severity: "danger",
        summary: "Error",
        detail: `failed to load tags (status ${err.response.status}): ${err.response.data}`,
        life: 10000,
      });
    })
    .finally(() => {
      loading.value = false;
    });
}

function reset() {
  pagination.page = 1;
  refresh();
}

function renderExpand(row: Tag) {
  return h("pre", null, JSON.stringify(row, null, 2));
}
</script>

<template>
  <Toast />
  <Toolbar>
    <template #center>
      <InputGroup>
        <XSearch v-model:value="params.search" @change="reset" />
        <Button @click="refresh" icon="pi pi-refresh" severity="secondary" />
      </InputGroup>
    </template>
    <template #end>
      <InputGroup>
        <Button
          label="Create"
          icon="pi pi-plus"
          @click="showNewTagModal = true"
        />
        <XTagDialog v-model:show="showNewTagModal" />
        <Button
          label="Delete"
          severity="danger"
          icon="pi pi-trash"
          :badge="selectedItems.length.toString()"
        />
      </InputGroup>
    </template>
  </Toolbar>
  <DataTable
    v-model:selection="selectedItems"
    :value="items"
    size="small"
    stripedRows
    showGridlines
  >
    <Column selectionMode="multiple" />
    <XColorColumn key="color" fieldName="color" title="Color" />
    <Column field="name" header="Name" />
    <Column field="description" header="description" />
  </DataTable>

  <XPagination
    v-model:page="page"
    v-model:size="size"
    :total="total"
    :more="more"
    @update="refresh"
  />
</template>
