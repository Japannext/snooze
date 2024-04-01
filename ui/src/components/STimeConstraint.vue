<script setup lang="ts">
import {
  NTag, NIcon,
  NSpace,
} from 'naive-ui'

import { Sunny, Calendar, CalendarNumber, ArrowForward, Time } from '@vicons/ionicons5'

interface Props {
  value: object,
}

const props = defineProps<Props>()

function printTime(datetime: string): string {
  if (datetime === null) {
    return 'null'
  }
  const time = datetime.split('T').pop()
  const timeSplit = time.split(':')
  return `${timeSplit[0]}:${timeSplit[1]}`
}

function printDate(datetime: string): string {
  return datetime.split('T')[0]
}

</script>

<template>
  <n-space v-if="value.kind == 'daily'" align="center" :size="0">
    <n-tag type="warning" round :bordered="false">
      <template #icon>
        <n-icon :component="Sunny" />
      </template>
      {{ printTime(value.start.time) }}
    </n-tag>
    <n-tag type="warning" round :bordered="false">
      <template #icon>
        <n-icon :component="ArrowForward" />
      </template>
      {{ printTime(value.end.time) }}
    </n-tag>
  </n-space>
  <n-space v-else-if="value.kind == 'weekly'" align="center" :size="0">
    <n-tag type="info" round :bordered="false">
      <template #icon>
        <n-icon :component="Calendar" />
      </template>
      <b>{{ value.start.weekday }}</b>
      {{ printTime(value.start.time) }}
    </n-tag>
    <n-tag type="info" round :bordered="false">
      <template #icon><n-icon :component="ArrowForward" /></template>
      <b>{{ value.end.weekday }}</b>
      {{ printTime(value.end.time) }}
    </n-tag>
  </n-space>
  <n-space v-else-if="value.kind == 'absolute'" align="center" :size="0">
    <n-tag type="error" round :bordered="false">
      <template #icon>
        <n-icon :component="CalendarNumber" />
      </template>
      <b>{{ printDate(value.start.datetime) }}</b>
      {{ printTime(value.start.datetime) }}
    </n-tag>
    <n-tag type="error" round :bordered="false">
      <template #icon>
        <n-icon :component="ArrowForward" />
      </template>
      <b>{{ printDate(value.end.datetime) }}</b>
      {{ printTime(value.end.datetime) }}
    </n-tag>
  </n-space>
  <n-tag v-else-if="value.kind == 'always'">
    <template #icon>
      <n-icon :component="Time" />
    </template>
    Always
  </n-tag>
</template>
