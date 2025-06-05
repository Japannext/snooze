<script setup lang="ts">
import { defineProps, withDefaults } from "vue";
import { useTimeRelative } from "@/stores";
import { DateTime } from "luxon";

const timeRelative = useTimeRelative();

interface Props {
  title?: string;
  fieldName: string;
}

withDefaults(defineProps<Props>(), {
  title: "Timestamp",
  fieldName: "timestamp",
});

function relativeTime(ts: number): string {
  if (ts) {
    return DateTime.fromMillis(ts, { zone: "system" }).toRelative();
  }
  return "undefined";
}

function absoluteTime(ts: number): string {
  if (ts) {
    return DateTime.fromMillis(ts, { zone: "system" }).toFormat(
      "yyyy-MM-dd HH:mm:ss",
    );
  }
  return "undefined";
}
</script>

<template>
  <Column :field="fieldName" style="width: 160px">
    <template #header>
      <ToggleButton
        v-model="timeRelative"
        size="small"
        :onLabel="title"
        onIcon="pi pi-hourglass"
        :offLabel="title"
        offIcon="pi pi-globe"
      />
    </template>
    <template #body="{ data }">
      <span v-if="timeRelative">
        {{ relativeTime(data[fieldName]) }}
      </span>
      <span v-else>
        {{ absoluteTime(data[fieldName]) }}
      </span>
    </template>
  </Column>
</template>
