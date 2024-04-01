<script setup lang="ts">
import { h, computed, defineEmits, onMounted } from 'vue'
import {
  NGrid, NGi,
  NInput, NInputGroup,
  NPopselect, NSelect,
  NButton, NIcon,
  NSpace, NList, NListItem, NEmpty,
} from 'naive-ui'
import { Add, Trash, Refresh } from '@vicons/ionicons5'
import type { SelectOption } from 'naive-ui'

import SConditionOperator from '@/components/inputs/SConditionOperator.vue'

interface Props {
  value: object,
  root: boolean,
  index: number,
  size: string,
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: "update:value", value: object):void,
  (e: "delete", value: number):void,
}>()

const dataValue = computed({
  get() {
    console.log(`child[${props.index}] dataValue.get(): ${JSON.stringify(props.value)}`)
    return props.value
  },
  set(v) { emit("update:value", v) },
})

const operationType = computed(() => {
  if (dataValue.value === undefined) {
    return 'none'
  }
  switch(dataValue.value.kind) {
    case 'and': case 'or':
      return 'logic'
    case 'not':
      return 'not'
    case 'always_true':
      return 'always_true'
    case undefined: case null:
      return 'none'
    default:
      return 'operator'
  }
})

onMounted(() => {
  if (dataValue.value === undefined || dataValue.value.kind === undefined) {
    emptyCondition()
  }
})

const logicOptions: SelectOption[] = [
  {value: 'and', label: "and"},
  {value: 'or', label: "or"},
  {value: 'not', label: "not"},
]

function defaultCondition(): object {
  return {kind: '=', field: ["resource"], value: ""}
}

function addCondition() {
  dataValue.value = {
    kind: 'and',
    conditions: [
      dataValue.value,
      defaultCondition(),
    ]
  }
}

function escalateDelete() {
  if (props.root) {
    console.log("Emptying condition...")
    emptyCondition()
    console.log(`dataValue = ${JSON.stringify(dataValue.value)}`)
  } else {
    emit("delete", props.index)
  }
}

function deleteNot() {
  dataValue.value = dataValue.value.condition
}

function deleteNotHandler(_idx: number) {
  dataValue.value.condition = defaultCondition()
}

function deleteLogicHandler(idx: number) {
  // Remove the correct condition
  dataValue.value.conditions.splice(idx, 1)

  // Remove the AND/OR if there is less than 2 arguments
  if (operationType.value == 'logic' && dataValue.value.conditions.length < 2) {
    dataValue.value = dataValue.value.conditions[0]
  }
}

function setDefaultCondition() {
  dataValue.value = {
    kind: '=',
    field: ['resource'],
    value: '',
  }
}

function emptyCondition() {
  dataValue.value = {kind: 'always_true'}
}

</script>

<template>
  <template v-if="['logic', 'not'].includes(operationType)">
    <n-grid :y-gap="5" :cols="100">
      <n-gi :span="100">
        <n-input-group>
          <n-select
            v-model:value="dataValue.kind"
            style="width: 70px;"
            :size="size"
            :options="logicOptions"
          />
          <n-button
            v-if="dataValue.kind == 'not'"
            secondary
            type="error"
            :size="size"
            @click="deleteNot"
          >
            <n-icon :component="Trash" />
          </n-button>
        </n-input-group>
      </n-gi>
      <template v-if="operationType == 'logic'">
        <n-gi
          v-for="(condition, idx) in dataValue.conditions" :key="`logic-${idx}`"
          :span="98" :offset="2"
        >
          <s-condition-child
            v-model:value="dataValue.conditions[idx]"
            :index="idx"
            :root="false"
            :size="size"
            @delete="deleteLogicHandler"
          />
        </n-gi>
      </template>
      <n-gi
        v-if="operationType == 'not'"
        :span="98" :offset="2"
      >
        <s-condition-child
          v-model:value="dataValue.condition"
          :index="0"
          :root="false"
          :size="size"
          @delete="deleteNotHandler"
        />
      </n-gi>
    </n-grid>
  </template>

  <template v-else-if="operationType == 'operator'">
    <n-input-group>
      <s-condition-operator v-model:value="dataValue" :size="size" />
      <n-button
        secondary
        type="error"
        :size="size"
        @click="escalateDelete"
      >
        <template #icon><n-icon :component="Trash" /></template>
      </n-button>
      <n-button
        secondary
        :size="size"
        @click="addCondition"
      >
        <template #icon><n-icon :component="Add" /></template>
      </n-button>
    </n-input-group>
  </template>

  <n-empty
    v-else-if="operationType == 'always_true'"
    description="Always true"
    :size="size"
  >
    <template #icon><n-icon :component="Refresh" /></template>
    <template #extra>
      <n-button @click="setDefaultCondition()">Add condition</n-button>
    </template>
  </n-empty>

  <n-empty
    v-else
    description="Invalid condition"
    :size="size"
  >
    <template #extra>
      <n-button @click="setDefaultCondition()">Add condition</n-button>
    </template>
  </n-empty>
</template>
