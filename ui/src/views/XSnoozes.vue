<script setup lang="ts">
import { h, ref, onMounted } from "vue";
import { useLoadingBar, useMessage } from "naive-ui";
import { usePagination } from "@/api";
import { useRouteQuery } from "@vueuse/router";
import { useToast } from "primevue/usetoast";
import { useDialog } from "primevue/usedialog";

import {
  getSnoozes,
  type Snooze,
  type GetSnoozesParams,
  type Filter,
} from "@/api";

const items = ref<Snooze[]>();
const selectedItems = ref<string[]>([]);
const xTimerange = ref();

const page = useRouteQuery<number>("page", 1, { transform: Number });
const size = useRouteQuery<number>("size", 20, { transform: Number });
const total = ref<number>(0);
const more = ref<boolean>(false);

const toast = useToast();
const dialog = useDialog();
const loading = ref<boolean>(false);
const showSnoozeDialog = ref<boolean>(false);
const showCancelDialog = ref<boolean>(false);

const params = ref<GetSnoozesParams>({
  search: "",
  pagination: {},
  filter: "active",
  timerange: {},
});

const filters: Filter[] = [
  { label: "Active", value: "active" },
  { label: "Upcoming", value: "upcoming" },
  { label: "Expired", value: "expired" },
  { label: "Cancelled", value: "cancelled" },
  { label: "All", value: "all" },
];

onMounted(() => {
  refresh();
});

function refresh() {
  console.log(`refresh()`);
  loading.value = true;
  params.value.timerange = xTimerange.value.getTime();
  params.value.pagination = {
    page: page.value,
    pageSize: size.value,
  };
  getSnoozes(params.value)
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
        detail: `failed to load snoozes: ${err}`,
      });
    })
    .finally(() => {
      loading.value = false;
    });
}

const delayMillis = 300;

function refreshWithDelay() {
  setTimeout(refresh, delayMillis);
}

function reset() {
  params.value.pagination.page = 1;
  refresh();
}
</script>

<template>
  <Toast />
  <Toolbar>
    <template #start>
      <XFilter
        v-model:value="params.filter"
        :filters="filters"
        @change="reset"
      />
    </template>
    <template #center>
      <InputGroup>
        <XTimeRange ref="xTimerange" @change="reset" />
        <XSearch v-model:value="params.search" @change="reset" />
        <Button @click="refresh" icon="pi pi-refresh" severity="secondary" />
      </InputGroup>
    </template>
    <template #end>
      <InputGroup>
        <Button
          label="New Snooze"
          severity="warn"
          icon="pi pi-plus"
          @click="showSnoozeDialog = true"
        />
        <XSnoozeDialog
          v-model:show="showSnoozeDialog"
          @success="refreshWithDelay"
        />
        <Button
          label="Cancel"
          severity="danger"
          icon="pi pi-times"
          :badge="selectedItems.length.toString()"
        />
        <XCancelDialog
          v-model:show="showCancelDialog"
          :ids="selectedItems"
          @success="refreshWithDelay"
        />
      </InputGroup>
    </template>
  </Toolbar>
  <XLoadingBar :loading="loading" />

  <DataTable
    v-model:selection="selectedItems"
    :value="items"
    size="small"
    stripedRows
    showGridlines
  >
    <Column selectionMode="multiple" />
    <XTimestampColumn key="from" fieldName="startAt" title="From" />
    <XTimestampColumn key="to" fieldName="endsAt" title="To" />
    <Column field="groups" header="Groups" />
    <XDurationColumn
      key="duration"
      fromField="startAt"
      toField="endsAt"
      title="Duration"
    />
    <Column field="username" header="Username" />
    <Column field="tags" header="Tags" />
    <Column field="reason" header="Reason" />
  </DataTable>

  <XPagination
    v-model:page="page"
    v-model:size="size"
    :total="total"
    :more="more"
    @update="refresh"
  />
</template>
