<script setup lang="ts">

import { ref, h, defineEmits, withDefaults, computed } from 'vue'
import {
  NInput, NInputNumber,
  NInputGroupLabel,
  NPopselect, NSelect,
  NButton, NIcon,
} from 'naive-ui'
import type { SelectOption, SelectGroupOption } from 'naive-ui'
import { Text, Ban, Checkmark } from '@vicons/ionicons5'
import { Boolean as BooleanIcon } from '@vicons/carbon'
import { NumberOutlined, QuestionCircleOutlined } from '@vicons/antd'
import { DecimalArrowRight20Filled } from '@vicons/fluent'
import { useVModel } from '@vueuse/core'

import type { AnyType } from '@/api'

interface Props {
  value: AnyType,
  size: string,
}

const props = withDefaults(defineProps<Props>(), {
  value: {kind: "string", data: ""},
  size: "medium",
})

const emit = defineEmits<{
  (e: "update:value", value: AnyType): void
}>()

const dataValue = useVModel(props, 'value', emit)

const selectOptions = [
  {label: "String", value: "string"},
  {label: "Integer", value: "integer"},
  {label: "Float", value: "float"},
  {label: "Boolean", value: "boolean"},
  {label: "Null", value: "null"},
]

const booleanOptions = [
  {label: "True", value: true, style: {color: "green"}},
  {label: "False", value: false, style: {color: "red"}},
]

function getDefault(kind: string) {
  switch(kind) {
    case "string":
      return ""
    case "integer": case "float":
      return 0
    case "boolean":
      return true
    case "null":
      return null
  }
}

function updateOnSelect(kind: string) {
  dataValue.value.kind = kind
  dataValue.value.data = getDefault(kind)
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
    case "integer":
      return h('span', {}, [
        h('NIcon', {component: NumberOutlined}),
        h('span', "Integer"),
      ])
    case "float":
      return h('span', {}, [
        h('NIcon', {component: DecimalArrowRight20Filled}),
        h('span', "Float"),
      ])
    case "boolean":
      return h('span', {}, [
        h('NIcon', {component: DecimalArrowRight20Filled}),
        h('span', "Boolean"),
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
    v-if="dataValue.kind == 'string'"
    v-model:value="dataValue.data"
    placeholder="String"
    :size="size"
  />
  <n-input-number
    v-else-if="dataValue.kind == 'integer'"
    v-model:value="dataValue.data"
    placeholder="Integer"
    step="1"
    :size="size"
  />
  <n-input-number
    v-else-if="dataValue.kind == 'float'"
    v-model:value="dataValue.data"
    placeholder="Float"
    step="0.1"
    :size="size"
  />
  <n-select
    v-else-if="dataValue.kind == 'boolean'"
    v-model:value="dataValue.data"
    style="width: 80px;"
    :options="booleanOptions"
    :size="size"
  />
  <n-input-group-label
    v-else-if="dataValue.kind == 'null'"
    placeholder="Null"
    :size="size"
  >
    Null
  </n-input-group-label>
  <n-input-group-label
    v-else
    :size="size"
  >
    {{ dataValue.data }}
  </n-input-group-label>

  <n-popselect
    v-model:value="dataValue.kind"
    :render-label="renderLabel"
    :options="selectOptions"
    :size="size"
    @update:value="updateOnSelect"
  >
    <n-button v-if="dataValue.kind == 'string'" :size="size"><n-icon :component="Text" /></n-button>
    <n-button v-else-if="dataValue.kind == 'integer'" :size="size"><n-icon :component="NumberOutlined" /></n-button>
    <n-button v-else-if="dataValue.kind == 'float'" :size="size"><n-icon :component="DecimalArrowRight20Filled" /></n-button>
    <n-button v-else-if="dataValue.kind == 'boolean'" :size="size"><n-icon :component="BooleanIcon" /></n-button>
    <n-button v-else-if="dataValue.kind == 'null'" :size="size"><n-icon :component="Ban" /></n-button>
    <n-button v-else :size="size"><n-icon :component="QuestionCircleOutined" /></n-button>
  </n-popselect>
</template>
