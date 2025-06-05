<script setup lang="ts">
import { defineProps } from "vue";
import { Duration, type DurationLikeObject, Interval, DateTime } from "luxon";

const props = defineProps<{
  fromField: string;
  toField: string;
  title: string;
}>();

type DurationUnit = keyof DurationLikeObject;

const day = Duration.fromObject({ days: 1 }).toMillis();
const hour = Duration.fromObject({ hours: 1 }).toMillis();
// const minute = Duration.fromObject({minutes: 1}).toMillis()
const units: DurationUnit[] = [
  "years",
  "months",
  "days",
  "hours",
  "minutes",
  "seconds",
  "milliseconds",
];

function toHuman(
  dur: Duration,
  smallestUnit: DurationUnit = "seconds",
): string {
  const smallestIdx = units.indexOf(smallestUnit);
  const entries = Object.entries(
    dur
      .shiftTo(...units)
      .normalize()
      .toObject(),
  ).filter(
    ([_unit, amount]: [string, number], idx) =>
      amount > 0 && idx <= smallestIdx,
  );
  const dur2 = Duration.fromObject(
    entries.length === 0 ? { [smallestUnit]: 0 } : Object.fromEntries(entries),
  );
  return dur2.toHuman();
}

// Get the severity color associated with the duration.
function getSeverity(d: number): string {
  const duration = Duration.fromMillis(d);
  switch (true) {
    case duration <= 1 * hour:
      return "default";
    case duration <= 6 * hour:
      return "success";
    case duration <= 1 * day:
      return "info";
    case duration <= 7 * day:
      return "warning";
    default:
      return "error";
  }
}

function toDuration(from: number, to: number): string {
  const start = DateTime.fromMillis(from);
  const end = DateTime.fromMillis(to);
  if (start && end) {
    const interval = Interval.fromDateTimes(start, end);
    return toHuman(interval.toDuration());
  }
  return "unknown";
}
</script>

<template>
  <Column :field="fieldName" :header="title">
    <template #body="{ data }">
      <Tag :value="toDuration(data[fromField], data[toField])" />
    </template>
  </Column>
</template>
