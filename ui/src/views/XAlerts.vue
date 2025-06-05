<script setup lang="ts">
import { h, ref, onMounted, Component } from "vue";
import { useRouteQuery } from "@vueuse/router";
import { usePagination } from "@/api";

import {
  getAlerts,
  type GetAlertsParams,
  type Alert,
  type Filter,
} from "@/api";

const selected = ref<Alert | null>(null);
const showDetails = ref<boolean>(false);

// Pagination variables
const page = useRouteQuery<number>("page", 1, { transform: Number });
const size = useRouteQuery<number>("size", 20, { transform: Number });
const total = ref<number>(0);
const more = ref<boolean>(false);

// Utils
const loading = ref<boolean>(false);

const items = ref<Array<Alert>>([]);
const selectedItems = ref<Array<string>>([]);
const xTimerange = ref();

const showAckDialog = ref<boolean>(false);
const showEscalateDialog = ref<boolean>(false);

const pagination = usePagination(refresh);
const params = ref<GetAlertsParams>({
  search: "",
  filter: "active",
  pagination: {},
  timerange: {},
});

const filters: Filter[] = [
  { label: "Active", value: "active" },
  { label: "History", value: "history" },
];

function refresh() {
  loading.value = true;
  params.value.timerange = xTimerange.value.getTime();
  params.value.pagination = {
    page: page.value,
    pageSize: size.value,
  };
  getAlerts(params.value)
    .then((list) => {
      items.value = list.items;
      total.value = list.total;
      more.value = list.more;
    })
    .catch((err) => {
      items.value = [];
      message.error(`failed to load alerts: ${err}`);
    })
    .finally(() => {
      loading.value = false;
    });
}

onMounted(() => {
  refresh();
});

function unselect() {
  selectedItems.value = [];
}

function select(item: Alert) {
  selected.value = item;
  showDetails.value = true;
}

function reset() {
  params.value.pagination.page = 1;
  refresh();
}

function renderExpand(row: Alert) {
  return h("pre", null, JSON.stringify(row, null, 2));
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
          label="Ack"
          icon="pi pi-check"
          @click="showAckDialog = true"
          :disabled="selectedItems.length == 0"
          :badge="selectedItems.length.toString()"
        />
        <XAckDialog
          v-model:show="showAckDialog"
          :ids="selectedItems"
          @success="unselect"
        />
        <Button
          label="Escalate"
          severity="warn"
          icon="pi pi-arrow-circle-up"
          :disabled="selectedItems.length == 0"
          :badge="selectedItems.length.toString()"
        />
      </InputGroup>
    </template>
  </Toolbar>
  <div style="height: 2px">
    <ProgressBar v-if="loading" mode="indeterminate" style="height: 2px" />
  </div>

  <DataTable
    v-model:selection="selectedItems"
    :value="items"
    size="small"
    stripedRows
    showGridlines
  >
    <Column selectionMode="multiple" />
    <XTimestampColumn key="timestamp" fieldName="startsAt" />
    <XAlertAttributeColumn key="attributes" />
    <Column field="summary" header="Summary" />
  </DataTable>

  <XPagination
    v-model:page="page"
    v-model:size="size"
    :total="total"
    :more="more"
    @update="refresh"
  />
</template>
