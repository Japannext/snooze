<script setup lang="ts">
import { NIcon, NButton, NPopover, NCollapse, NCollapseItem, NDatePicker, NRadioGroup, NRadioButton, NTabs, NTabPane } from 'naive-ui'
import { ref, onMounted, defineModel, computed } from 'vue'
import type { PropType } from 'vue'
import { ClockRegular } from '@vicons/fa'
import { DateTime } from 'luxon'

const rangeMillis = defineModel<[number, number]>("rangeMillis")

const rangeText = ref<string>(null)
const emit = defineEmits(['refresh'])
const showPopover = ref<Boolean>(false)

const preset = ref<string>("last_24_hours")
const day = (24 * 3600 * 1000)

function setPreset(preset: string) {
  console.log(`setPreset(${preset})`)
  var now = Date.now()
  switch(preset) {
    case "last_24_hours":
      rangeMillis.value = [now - day, now]
      rangeText.value = "last 24 hours"
      break
    case "last_7_days":
      rangeMillis.value = [now - (7 * day), now]
      rangeText.value = "last 7 days"
      break
  }
  ok()
}

function setTimeRange(value: [number, number]) {
  console.log(`setTimeRange(${value})`)
  preset.value = null
  rangeMillis.value = value
  var fromDate = DateTime.fromMillis(value[0]).toISODate()
  var toDate = DateTime.fromMillis(value[1]).toISODate()
  rangeText.value = `${fromDate} -> ${toDate}`
  ok()
}

function ok() {
  console.log(`ok()`)
  showPopover.value = false
  emit('refresh')
}

function exit() {
  console.log(`exit()`)
  showPopover.value = false
}

onMounted(() => {
  setPreset("last_24_hours")
})
</script>

<template>
  <n-popover trigger="manual" :show="showPopover" @on-clickoutside="exit">
    <template #trigger>
      <n-button @click="showPopover = !showPopover">
        <template #icon><n-icon :component="ClockRegular" /></template>
        {{ rangeText }}
      </n-button>
    </template>
    <n-tabs>
      <n-tab-pane name="presets" tab="Presets">
        <n-radio-group v-model:value="preset">
          <n-radio-button value="last_24_hours" label="Last 24 hours" @click="setPreset('last_24_hours')" />
          <n-radio-button value="last_7_days" label="Last 7 days" @click="setPreset('last_7_days')" />
        </n-radio-group>
      </n-tab-pane>
      <n-tab-pane name="timerange" tab="Time range">
        <n-date-picker
          v-model:value="rangeMillis"
          type="daterange"
          clearable
          panel
          :on-confirm="setTimeRange"
        />
      </n-tab-pane>
    </n-tabs>
  </n-popover>
</template>
