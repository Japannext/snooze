<script setup lang="ts">

import { computed, withDefaults, defineEmits } from 'vue'
import { format, startOfToday, endOfTomorrow } from 'date-fns'

import {
  NDatePicker, NTimePicker,
  NInputGroup, NInputGroupLabel,
  NSelect,
  NRadio, NRadioGroup,
  NSpace, NEmpty,
  NIcon,
} from 'naive-ui'
import { Time } from '@vicons/ionicons5'
import type { SelectOption } from 'naive-ui'

interface Props {
  value: object,
  size: string,
}

const props = withDefaults(defineProps<Props>(), {
  size: "small",
})

const emit = defineEmits<{
  (e: "update:value", value: object): void,
}>()

const dataValue = computed({
  get() { return props.value },
  set(v) { emit("update:value", v) }
})

const datetimeValue = computed({
  get() { return [dataValue.value.start.datetime, dataValue.value.end.datetime] },
  set(v) {
    dataValue.value.start.datetime = v[0]
    dataValue.value.end.datetime = v[1]
  },
})

const DATE_FORMAT = "yyyy-MM-dd'T'HH:mm:ss.SSSSSSXX"

function getDefault(kind: string): object {
  switch(kind) {
    case "daily":
      return {'kind': 'daily', 'start': {time: "00:00:00"}, 'end': {time: "23:59:59"}}
    case "weekly":
      return {'kind': 'weekly', 'start': {weekday: "monday", time: "00:00:00"}, 'end': {weekday: "friday", time: "23:59:59"}}
    case "absolute":
      return {'kind': 'absolute',
        'start': {datetime: format(startOfToday(), DATE_FORMAT)},
        'end': {datetime: format(endOfTomorrow(), DATE_FORMAT)},
      }
    case "always":
      return {'kind': 'always'}
  }
}

const dataKind = computed({
  get() { return dataValue.value.kind },
  set(v) { dataValue.value = getDefault(v) },
})

const weekdayOptions: SelectOption[] = [
  {value: "monday", label: "Monday"},
  {value: "tuesday", label: "Tuesday"},
  {value: "wednesday", label: "Wednesday"},
  {value: "thursday", label: "Thursday"},
  {value: "friday", label: "Friday"},
  {value: "saturday", label: "Saturday"},
  {value: "sunday", label: "Sunday"},
]

</script>

<template>
  <n-space vertical>
    <n-radio-group v-model:value="dataKind">
      <n-radio value="daily">Daily</n-radio>
      <n-radio value="weekly">Weekly</n-radio>
      <n-radio value="absolute">Absolute</n-radio>
      <n-radio value="always">Always</n-radio>
    </n-radio-group>

    <n-space v-if="dataValue.kind == 'daily'">
      <n-input-group>
        <n-input-group-label :size="size">From</n-input-group-label>
        <n-time-picker
          v-model:formatted-value="dataValue.start.time"
          style="min-width: 80px;"
          value-format="HH:mm:ss"
          :size="size"
          :actions="['now']"
        />
      </n-input-group>
      <n-input-group>
        <n-input-group-label :size="size">From</n-input-group-label>
        <n-time-picker
          v-model:formatted-value="dataValue.end.time"
          style="min-width: 80px;"
          value-format="HH:mm:ss"
          :size="size"
        />
      </n-input-group>
    </n-space>

    <n-space v-else-if="dataValue.kind == 'weekly'">
      <n-input-group>
        <n-input-group-label :size="size">From</n-input-group-label>
        <n-select
          v-model:value="dataValue.start.weekday"
          style="min-width: 120px; max-width: 120px;"
          :size="size"
          :options="weekdayOptions"
        />
        <n-time-picker
          v-model:formatted-value="dataValue.start.time"
          style="min-width: 80px;"
          value-format="HH:mm:ss"
          :size="size"
          :actions="['now']"
        />
      </n-input-group>
      <n-input-group>
        <n-input-group-label :size="size">To</n-input-group-label>
        <n-select
          v-model:value="dataValue.end.weekday"
          style="min-width: 120px; width: 120px;"
          :size="size"
          :options="weekdayOptions"
        />
        <n-time-picker
          v-model:formatted-value="dataValue.end.time"
          value-format="HH:mm:ss"
          style="min-width: 80px;"
          :size="size"
        />
      </n-input-group>
    </n-space>

    <n-space v-else-if="dataValue.kind == 'absolute'">
      <n-date-picker
        v-model:formatted-value="datetimeValue"
        :value-format="DATE_FORMAT"
        type="datetimerange"
        :size="size"
        clearable
      />
    </n-space>

    <n-empty
      v-else-if="dataValue.kind == 'always'"
      description="Always"
    >
      <template #icon><n-icon :component="Time" /></template>
    </n-empty>
  </n-space>
</template>
