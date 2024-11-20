<script setup lang="ts">
import { defineProps } from 'vue'
import { Duration } from 'luxon'

import { NTag } from 'naive-ui'

const props = defineProps<{
  duration: number,
}>()

const day = Duration.fromObject({days: 1}).toMillis()
const hour = Duration.fromObject({hours: 1}).toMillis()
// const minute = Duration.fromObject({minutes: 1}).toMillis()

function toHuman(dur: Duration, smallestUnit: string = "seconds"): string {
  const units = ["years", "months", "days", "hours", "minutes", "seconds", "milliseconds", ];
  const smallestIdx = units.indexOf(smallestUnit);
  const entries = Object.entries(
    dur.shiftTo(...units).normalize().toObject()
  ).filter(([_unit, amount], idx) => amount > 0 && idx <= smallestIdx);
  const dur2 = Duration.fromObject(
    entries.length === 0 ? { [smallestUnit]: 0 } : Object.fromEntries(entries)
  );
  return dur2.toHuman();
}

function getTagType(): string {
  switch (true) {
  case (props.duration <= 1 * hour):
    return ""
  case (props.duration <= 6 * hour):
    return "success"
  case (props.duration <= 1 * day):
    return "info"
  case (props.duration <= 7 * day):
    return "warning"
  default:
    return "error"
  }
}

</script>

<template>
  <n-tag round :type="getTagType()">
    {{ toHuman(Duration.fromMillis(duration)) }}
  </n-tag>
</template>
