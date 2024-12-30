<script setup lang="ts">
import { defineProps } from 'vue'
import { Duration, type DurationLikeObject } from 'luxon'

import { NTag } from 'naive-ui'
import { type NaiveColor } from '@/api'

const props = defineProps<{
  duration: number,
}>()

type DurationUnit = keyof DurationLikeObject

const day = Duration.fromObject({days: 1}).toMillis()
const hour = Duration.fromObject({hours: 1}).toMillis()
// const minute = Duration.fromObject({minutes: 1}).toMillis()
const units: DurationUnit[] = ["years", "months", "days", "hours", "minutes", "seconds", "milliseconds" ]

function toHuman(dur: Duration, smallestUnit: DurationUnit = "seconds"): string {
  const smallestIdx = units.indexOf(smallestUnit);
  const entries = Object.entries(
    dur.shiftTo(...units).normalize().toObject()
  ).filter(([_unit, amount]: [string, number], idx) => amount > 0 && idx <= smallestIdx);
  const dur2 = Duration.fromObject(
    entries.length === 0 ? { [smallestUnit]: 0 } : Object.fromEntries(entries)
  );
  return dur2.toHuman();
}

function getTagType(): NaiveColor {
  switch (true) {
  case (props.duration <= 1 * hour):
    return "default"
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
