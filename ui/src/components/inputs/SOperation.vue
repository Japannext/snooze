<script setup lang="ts">

import { computed, defineEmits } from 'vue'
import { NPopselect, NButton } from 'naive-ui'
import type { SelectOption } from 'naive-ui'

interface Props {
  value: string,
  size: string,
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: "update:value", value: string): void,
}>()

interface Option {
  text: string
}

const optionsMap: Map<string, Option> = new Map([
  ["=", {text: "="}],
  ["!=", {text: "≠"}],
  [">", {text: ">"}],
  ["<", {text: "<"}],
  [">=", {text: "≥"}],
  ["<=", {text: "≤"}],
  ["matches", {text: "matches"}],
  ["exists", {text: "exists?"}]
])

const options: SelectOption[] = [
  {value: "=", label: "="},
  {value: "!=", label: "≠"},
  {value: ">", label: ">"},
  {value: "<", label: "<"},
  {value: ">=", label: "≥"},
  {value: "<=", label: "≤"},
  {value: "matches", label: "matches"},
  {value: "exists", label: "exists?"},
]


const dataValue = computed({
  get() { return props.value },
  set(v) {
    emit("update:value", v)
  },
})

const option = computed(() => {
  return optionsMap.get(dataValue.value)
})

</script>

<template>
  <n-popselect
    v-model:value="dataValue"
    trigger="click"
    :options="options"
    :size="size"
  >
    <n-button :size="size">{{ option.text }}</n-button>
  </n-popselect>
</template>
