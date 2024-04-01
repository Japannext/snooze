<script setup lang="ts">
import { h, computed, defineEmits } from 'vue'
import {
  NGrid, NGi,
  NInput, NInputGroup,
  NPopselect,
} from 'naive-ui'
import type { SelectOption } from 'naive-ui'

import SMultiValue from '@/components/inputs/SMultiValue.vue'
import SOperation from '@/components/inputs/SOperation.vue'
import SField from '@/components/inputs/SField.vue'
import SAnyType from '@/components/inputs/SAnyType.vue'

interface Props {
  value: object,
  size: string,
}

const emit = defineEmits<{
  (e: "update:value", value: object):void,
}>()

const props = defineProps<Props>()

const dataValue = computed({
  get() { return props.value },
  set(v) { emit("update:value", v) },
})

const valueType = computed(() => {
  switch(props.value.kind) {
    case '=': case '!=':
      return "anytype"
    case 'matches':
      return "string"
    case '>': case '<': case '>=': case '<=':
      return 'comparable'
    case 'exists':
      return 'unary'
    default:
      console.log(`Error: unknown condition kind ${props.value.kind}`)
      return ''
  }
})

function getDefault(kind: string) {
  switch(kind) {
    case '=': case '!=':
    case '>': case '<': case '>=': case '<=':
      return {kind: "string", data: ""}
    case "matches":
      return ""
    case "exists":
      return null
    default:
      return null
  }
}

function updateOnSelect(kind: string) {
  dataValue.value.kind = kind
  dataValue.value.value = getDefault(kind)
}

</script>

<template>
  <s-field v-model:value="dataValue.field" :size="size" />
  <s-operation v-model:value="dataValue.kind" :size="size" @update:value="updateOnSelect" />
  <s-any-type
    v-if="valueType == 'anytype'"
    v-model:value="dataValue.value"
    :size="size"
  />
  <n-input
    v-else-if="valueType == 'string'"
    v-model:value="dataValue.value"
    :size="size"
  />
</template>
