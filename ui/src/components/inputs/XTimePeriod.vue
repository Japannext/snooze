<script setup lang="ts">
import { defineModel } from 'vue'
import { DateTime, Duration } from 'luxon'

import { NSpace, NInputGroup, NDatePicker, NButtonGroup, NButton } from 'naive-ui'

const start = defineModel<number>('start')
const end = defineModel<number>('end')

const day = Duration.fromObject({days: 1}).toMillis()
const hour = Duration.fromObject({hours: 1}).toMillis()
const minute = Duration.fromObject({minutes: 1}).toMillis()

function preset15min() {
  var now = DateTime.now().toMillis()
  start.value = now
  end.value = now + (15 * minute)
}

function preset1hour() {
  var now = DateTime.now().toMillis()
  start.value = now
  end.value = now + (1 * hour)
}

function preset6hours() {
  var now = DateTime.now().toMillis()
  start.value = now
  end.value = now + (6 * hour)
}

function preset1day() {
  var now = DateTime.now().toMillis()
  start.value = now
  end.value = now + (1 * day)
}

function preset7days() {
  var now = DateTime.now().toMillis()
  start.value = now
  end.value = now + (7 * day)
}
</script>

<template>
  <n-space justify="start" size="small" :wrap="false">
    <n-input-group>
      <n-date-picker v-model:value="start" type="datetime" placeholder="From" />
      <n-date-picker v-model:value="end" type="datetime" placeholder="To" />
    </n-input-group>
    <n-button-group>
      <n-button @click="preset15min">15 min</n-button>
      <n-button @click="preset1hour">1 hour</n-button>
      <n-button @click="preset6hours">6 hours</n-button>
      <n-button @click="preset1day">1 day</n-button>
      <n-button @click="preset7days">7 days</n-button>
    </n-button-group>
  </n-space>
</template>
