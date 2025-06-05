<script setup lang="ts">
import { pageSizes } from "@/api";
import type PageState from "primevue/paginator";

interface Props {
  total: number;
  more: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  total: 0,
  more: false,
});

const page = defineModel("page", { type: Number, default: 1 });
const size = defineModel("size", { type: Number, default: 20 });

const emit = defineEmits(["update"]);

function updatePage(event: PageState) {
  page.value = event.page + 1;
  size.value = event.rows;
  emit("update");
}

function currentPageReportTemplate(): string {
  if (props.more) {
    return "{totalRecords}+ records";
  }
  return "{totalRecords} records";
}
</script>

<template>
  <Paginator
    :first="(page - 1) * size"
    :rows="size"
    :totalRecords="props.total"
    :rowsPerPageOptions="pageSizes"
    @page="updatePage"
    template="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown CurrentPageReport"
    :currentPageReportTemplate="currentPageReportTemplate()"
  >
  </Paginator>
</template>
