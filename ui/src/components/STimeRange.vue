<script setup lang="ts">
import { NIcon, NButton, NPopover, NDatePicker, NRadioGroup, NRadioButton, NTabs, NTabPane } from 'naive-ui'
import { ref, onMounted, defineExpose } from 'vue'
import { ClockRegular } from '@vicons/fa'
import { DateTime } from 'luxon'
import { useRouteQuery } from '@vueuse/router'

const day = (24 * 3600 * 1000)

enum Preset {
  Last24Hours = "Last 24 hours",
  Last7Days = "Last 7 days"
}

const preset = useRouteQuery('preset', "")
const dateRange = ref<[number, number]>()
const start = useRouteQuery('start', 0, {transform: Number})
const end = useRouteQuery('end', 0, {transform: Number})

const showPopover = ref<Boolean>(false)
const displayText = ref<string>(Preset.Last24Hours)

const emit = defineEmits(['update'])
defineExpose({getTime})

onMounted(() => {
  if (preset.value == "" && start.value == 0 && end.value == 0) {
    setPreset(Preset.Last24Hours)
  }
})


function setPreset(newPreset: Preset) {
  preset.value = newPreset
  displayText.value = newPreset
  start.value = 0
  end.value = 0
  ok()
}

function setDateRange() {
  if (dateRange.value[0] > 0) {
    start.value = dateRange.value[0]
  } else {
    delete(start.value)
  }
  if (dateRange.value[1] > 0) {
    end.value = dateRange.value[1]
  } else {
    delete(end.value)
  }
  preset.value = ""
  ok()
}

// Method to call from outside to return the
// current timerange
function getTime(): [number, number] {
  if (preset.value) {
    var now = DateTime.now().toMillis()
    switch(preset.value) {
      case Preset.Last24Hours: {
        return [now - day, 0]
      }
      case Preset.Last7Days: {
        return [now - 7*day, 0]
      }
    }
  }
  if (start.value || end.value) {
    return [start.value, end.value]
  }
  return [null, null]
}

function ok() {
  showPopover.value = false
  emit('update')
}

function exit() {
  showPopover.value = false
}
</script>

<template>
  <n-popover trigger="manual" :show="showPopover" @on-clickoutside="exit">
    <template #trigger>
      <n-button @click="showPopover = !showPopover">
        <template #icon><n-icon :component="ClockRegular" /></template>
        {{ displayText }}
      </n-button>
    </template>
    <n-tabs>
      <n-tab-pane name="presets" tab="Presets">
        <n-radio-group v-model:value="preset">
          <n-radio-button :value="Preset.Last24Hours" label="Last 24 hours" @click="setPreset(Preset.Last24Hours)" />
          <n-radio-button :value="Preset.Last7Days" label="Last 7 days" @click="setPreset(Preset.Last7Days)" />
        </n-radio-group>
      </n-tab-pane>
      <n-tab-pane name="timerange" tab="Time range">
        <n-date-picker
          v-model:value="dateRange"
          type="daterange"
          clearable
          panel
          :on-confirm="setDateRange"
        />
      </n-tab-pane>
    </n-tabs>
  </n-popover>
</template>
