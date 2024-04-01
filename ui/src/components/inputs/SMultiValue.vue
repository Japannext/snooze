<script setup lang="ts">

import { ref, h, defineEmits, withDefaults, computed } from 'vue'

import {
  NInput, NInputNumber,
  NInputGroupLabel,
  NPopselect,
  NButton, NIcon,
} from 'naive-ui'

import type { SelectOption, SelectGroupOption } from 'naive-ui'

import { Text, Ban } from '@vicons/ionicons5'
import { NumberOutlined, QuestionCircleOutlined } from '@vicons/antd'

type MultiValue = string|number|null

function detectKind(): string{
  if (props.value === null) {
    return "null"
  }
  if (props.value !== undefined) {
    return (typeof props.value)
  } else {
    return "string"
  }
}

const kind = ref<string>(detectKind())

interface Props {
  value: string|number,
  size: string,
}

const props = withDefaults(defineProps<Props>(), {
  value: "",
  size: "medium",
})

const emit = defineEmits<{
  (e: "update:value", value: MultiValue): void
}>()

const dataValue = computed({
  get() { return props.value },
  set(v) { emit('update:value', v) },
})

const selectOptions = [
  {label: "String", value: "string"},
  {label: "Number", value: "number"},
  {label: "Null", value: "null"},
]

function updateOnSelect(value: string) {
  switch(kind.value) {
    case "string":
      if (Number.isInteger(dataValue.value)) {
        dataValue.value = dataValue.value.toString()
      } else {
        dataValue.value = ""
      }
      break
    case "number":
      if (Number.isInteger(dataValue.value)) {
        dataValue.value = Number(dataValue.value)
      } else {
        dataValue.value = null
      }
      break
    case "null":
      dataValue.value = null
      break
    default:
      dataValue.value = null
  }
}

/** Render a custom label for n-popselect
 * @param option The option to render
 * @return {VNodeChildren}
 */
function renderLabel (option: SelectOption | SelectGroupOption) {
  switch(option.value) {
    case "string":
      return h('span', {}, [
        h('NIcon', {component: Text}),
        h('span', "String"),
      ])
    case "number":
      return h('span', {}, [
        h('NIcon', {component: NumberOutlined}),
        h('span', "Integer"),
      ])
    case "null":
      return h('span', {}, [
        h('NIcon', {component: Ban}),
        h('span', "Null"),
      ])
    default:
      return h('span', {}, [
          h('NIcon', {component: QuestionCircleOutlined}),
          h('span', option.value),
        ])
  }
}

</script>

<template>
  <n-input
    v-if="kind == 'string'"
    v-model:value="dataValue"
    placeholder="String"
    :size="size"
  />
  <n-input-number
    v-else-if="kind == 'number'"
    v-model:value="dataValue"
    placeholder="Integer"
    :size="size"
  />
  <n-input-group-label
    v-else-if="kind == 'null'"
    placeholder="Null"
    :size="size"
  >
    Null
  </n-input-group-label>
  <n-input
    v-else
    v-model:value="dataValue"
    :placeholder="`Unknown type (${kind})`"
    :size="size"
  />

  <n-popselect
    v-model:value="kind"
    :render-label="renderLabel"
    :options="selectOptions"
    :size="size"
    @update:value="updateOnSelect"
  >
    <n-button v-if="kind == 'string'" :size="size"><n-icon :component="Text" /></n-button>
    <n-button v-else-if="kind == 'number'" :size="size"><n-icon :component="NumberOutlined" /></n-button>
    <n-button v-else-if="kind == 'null'" :size="size"><n-icon :component="Ban" /></n-button>
    <n-button v-else :size="size"><n-icon :component="QuestionCircleOutined" /></n-button>
  </n-popselect>
</template>
