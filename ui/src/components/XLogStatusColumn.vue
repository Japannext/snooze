<script setup lang="ts">
import { defineProps, computed } from "vue";
import type { LogStatus, LogStatusKind, NaiveColor } from "@/api";

const props = defineProps<{
  fieldName: string;
}>();

interface TagKinds {
  [key: number]: TagKind;
}

type TagKind = {
  type: NaiveColor;
  value: string;
};

const KINDS: TagKinds = {
  0: { type: "secondary", value: "active" },
  1: { type: "warn", value: "snoozed" },
  2: { type: "info", value: "silenced" },
  3: { type: "danger", value: "ratelimited" }, // SHOULD NOT HAPPEN
  4: { type: "danger", value: "dropped" }, // SHOULD NOT HAPPEN
  5: { type: "danger", value: "activecheck" }, // SHOULD NOT HAPPEN
  6: { type: "success", value: "acked" },
};

const defaultValue = "???";
const defaultSeverity = "secondary";

function getValue(status: LogStatus): string {
  if (status) {
    var k = KINDS[status.kind];
    if (k) {
      return k.value;
    }
    return defaultValue;
  }
  return defaultValue;
}

function getSeverity(status: LogStatus): string {
  if (status) {
    var k = KINDS[status.kind];
    if (k) {
      return k.type;
    }
    return defaultSeverity;
  }
  return defaultSeverity;
}
</script>

<template>
  <Column :field="fieldName" header="Status">
    <template #body="{ data }">
      <Tag
        :value="getValue(data.status)"
        :severity="getSeverity(data.status)"
      />
    </template>
  </Column>
</template>
