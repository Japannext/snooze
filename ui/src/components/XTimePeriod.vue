<script setup lang="ts">
import { ref, computed } from "vue";
import { DateTime, Duration } from "luxon";

const start = defineModel<number>("start", { type: Number, default: 0 });
const end = defineModel<number>("end", { type: Number, default: 0 });

const startValue = computed({
  get(): Date {
    return DateTime.fromMillis(start.value).toJSDate();
  },
  set(newValue: Date) {
    start.value = DateTime.fromJSDate(newValue).toMillis();
  },
});

const endValue = computed({
  get(): Date {
    return DateTime.fromMillis(end.value).toJSDate();
  },
  set(newValue: Date) {
    end.value = DateTime.fromJSDate(newValue).toMillis();
  },
});

const day = Duration.fromObject({ days: 1 }).toMillis();
const hour = Duration.fromObject({ hours: 1 }).toMillis();
const minute = Duration.fromObject({ minutes: 1 }).toMillis();

const preset = ref<string>("");
const presets = [
  { label: "15 min", value: "15min" },
  { label: "1 hour", value: "1hour" },
  { label: "6 hours", value: "6hours" },
  { label: "1 day", value: "1day" },
  { label: "7 days", value: "7days" },
];

function setDates(duration: Duration) {
  console.info(`setDates(${JSON.stringify(duration)})`);
  const now = DateTime.now();
  startValue.value = now.toJSDate();
  endValue.value = now.plus(duration).toJSDate();
}

function updatePreset(preset: string) {
  console.info(`updatePreset(${JSON.stringify(preset)})`);
  const now = DateTime.now();
  switch (preset) {
    case "15min":
      setDates(15 * minute);
      break;
    case "1hour":
      setDates(1 * hour);
      break;
    case "6hours":
      setDates(6 * hour);
      break;
    case "1day":
      setDates(1 * day);
      break;
    case "7days":
      setDates(7 * day);
      break;
  }
}

const datepickerPtOptions = {
  pcInputText: {
    props: {
      style: "width: 50px;",
    },
  },
};
</script>

<template>
  <div>
    <InputGroup>
      <IftaLabel>
        <DatePicker
          inputId="startDate"
          v-model="startValue"
          dateFormat="yy-mm-dd"
          hourFormat="24"
          :pt="datepickerPtOptions"
          showTime
        />
        <label for="startDate">Start</label>
      </IftaLabel>
      <IftaLabel>
        <DatePicker
          inputId="endDate"
          v-model="endValue"
          dateFormat="yy-mm-dd"
          hourFormat="24"
          showTime
        />
        <label for="endDate">End</label>
      </IftaLabel>
      <SelectButton
        :modelValue="preset"
        :options="presets"
        optionLabel="label"
        optionValue="value"
        @update:modelValue="updatePreset"
      />
    </InputGroup>
  </div>
</template>
