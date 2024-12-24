<script setup lang="ts">
import { defineProps, ref, computed } from 'vue'
import { DateTime, Duration } from 'luxon'

import { NSpace, NTag, NIcon } from 'naive-ui'

import { ArrowForwardSharp } from '@/icons'

const props = defineProps<{
  start: number,
  end: number,
  cancelled: boolean,
}>()

const startDate = computed(() => DateTime.fromMillis(props.start))
const endDate = computed(() => DateTime.fromMillis(props.end))

const kind = computed(() => {
  var now = DateTime.now()
  switch(true) {
  case (props.cancelled):
    return "cancelled"
  case (startDate.value <= now && now <= endDate.value):
    return "active"
  case (endDate.value < now):
    return "expired"
  case (now < startDate.value):
    return "upcoming"
  default:
  return "unknown"
  }
})

const SIZE = "20px"

const FORMAT = "yyyy-MM-dd HH:mm:ss"

</script>

<template>
  <n-space v-if="kind == 'cancelled'" align="center">
    <n-tag type="error">{{ startDate.toFormat(FORMAT) }}</n-tag>
    <n-tag checkable :checked="false">
      <template #icon><n-icon :component="ArrowForwardSharp" /></template>
    </n-tag>
    <n-tag type="error">{{ endDate.toFormat(FORMAT) }}</n-tag>
  </n-space>

  <n-space v-if="kind == 'expired'" align="center">
    <n-tag type="warning">{{ startDate.toFormat(FORMAT) }}</n-tag>
    <n-tag checkable :checked="false">
      <template #icon><n-icon :component="ArrowForwardSharp" /></template>
    </n-tag>
    <n-tag type="warning">{{ endDate.toFormat(FORMAT) }}</n-tag>
  </n-space>

  <n-space v-if="kind == 'upcoming'" align="center">
    <n-tag type="info">{{ startDate.toFormat(FORMAT) }}</n-tag>
    <n-icon :size="SIZE" :component="ArrowForwardSharp" />
  </n-space>

  <n-space v-if="kind == 'active'" align="center" size="small">
    <n-icon :component="ArrowForwardSharp" />
    <n-tag type="success"><b>{{ endDate.toFormat(FORMAT) }}</b></n-tag>
  </n-space>
</template>
