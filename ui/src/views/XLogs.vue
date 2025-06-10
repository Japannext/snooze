<script setup lang="ts">
import { h, ref, onMounted } from "vue";
import { useRouteQuery } from "@vueuse/router";
import { useToast } from "primevue/usetoast";
import { usePagination } from "@/api";
import {
  getLogs,
  pageSizes,
  type Log,
  type GetLogsParams,
  type Filter,
} from "@/api";

const loading = ref<boolean>(false);
const items = ref<Array<Log>>();
const selectedItems = ref<Array<string>>([]);
const expandedRows = ref({});
const xTimerange = ref();

// const message = useMessage();
const toast = useToast();
const showAckDialog = ref<boolean>(false);
const showEscalateDialog = ref<boolean>(false);

const page = useRouteQuery<number>("page", 1, { transform: Number });
const size = useRouteQuery<number>("size", 20, { transform: Number });
const total = ref<number>(0);
const more = ref<boolean>(false);

const pagination = usePagination(refresh);
const params = ref<GetLogsParams>({
  search: "",
  filter: "active",
  pagination: {},
  timerange: {},
});

const filters: Filter[] = [
  { label: "Active", value: "active" },
  { label: "Snoozed", value: "snoozed" },
  { label: "Acked", value: "acked" },
  { label: "All", value: "all" },
];

onMounted(() => {
  refresh();
});

function refresh() {
  loading.value = true;
  params.value.timerange = xTimerange.value.getTime();
  params.value.pagination = {
    page: page.value,
    pageSize: size.value,
  };
  getLogs(params.value)
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
        detail: `failed to load logs (status ${err.response.status}): ${err.response.data}`,
        life: 10000,
      });
    })
    .finally(() => {
      loading.value = false;
    });
}

function unselect() {
  selectedItems.value = [];
  refresh();
}

function renderExpand(row: Log) {
  return h("pre", null, JSON.stringify(row, null, 2));
}

// Reset the pagination and refresh.
function reset() {
  params.value.pagination.page = 1;
  refresh();
}

// Event received when pagination is changed.
function updatePage(event: PageState) {
  console.info(`updatePage(${JSON.stringify(event)})`);
  page.value = event.page + 1;
  size.value = event.rows;
  refresh();
}
</script>

<template>
  <div>
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
      v-model:expandedRows="expandedRows"
      dataKey="_id"
      :value="items"
      size="small"
      stripedRows
      showGridlines
    >
      <Column selectionMode="multiple" />
      <Column expander />
      <XLogStatusColumn key="status" fieldName="status" />
      <XTimestampColumn key="timestamp" fieldName="displayTime" />
      <XLogAttributeColumn key="attributes" />
      <XMessageColumn key="message" />

      <template #expansion="slotProps">
        <pre>{{ slotProps.data }}</pre>
      </template>
    </DataTable>

    <XPagination
      v-model:page="page"
      v-model:size="size"
      :total="total"
      :more="more"
      @update="refresh"
    />
  </div>
</template>
