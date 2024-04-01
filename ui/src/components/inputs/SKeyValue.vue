<script setup lang="ts">

import {
  h, ref, defineEmits,
  VNodeChild,
} from 'vue'

import {
  NCard, NEmpty,
  NAutoComplete,
  NInput, NInputGroup, NInputGroupLabel,
  NButton, NIcon, NTag,
  SelectOption, SelectGroupOption,
  NGrid, NGi,
} from 'naive-ui'

import { Add, Trash } from '@vicons/ionicons5'

import SMultiValue from '@/components/inputs/SMultiValue.vue'

const inputKey = ref(null)

type SelectRenderer = (option: SelectOption | SelectGroupOption) => VNodeChild

const size = "small"

interface Props {
  title: string,
  value: Map<string, string|number>,
  autoComplete: boolean,
  acOptions: string[],
  acRender: SelectRenderer,
}

const props = withDefaults(defineProps<Props>(), {
  title: "Name",
  value: () => new Map<string, string|number>(),
  autoComplete: false,
  acOptions: () => [],
  acRender: undefined,
})

const emit = defineEmits<{
  (e: "update:value", value: Map<string, string|number>): void
}>()


var dataValue = ref<Map<string, string|number>>(props.value)

const newItemKey = ref<string>("")
const newItemValue = ref<string|number>("")
const newItemType = ref<string>("string")

/** Reset the values of the new item
 */
function resetNewItem() {
  newItemKey.value = ""
  newItemValue.value = ""
  newItemType.value = "string"
}

/** Add an element to the key-value object
 */
function append() {
  dataValue.value.set(newItemKey.value, newItemValue.value)
  emit('update:value', dataValue.value)
  resetNewItem()
  // Focus back on the first input after appending
  inputKey.value.focus()
}

/** Update an element of the key-value object
 */
function update(key: string, val: string|number) {
  dataValue.value.set(key, val)
  emit('update:value', dataValue.value)
}

/** Remove an element from the key-value object
 */
function remove(key: string) {
  dataValue.value.delete(key)
  emit('update:value', dataValue.value)
}

function selectColor(label: string): string {
  var base = label.split('.')[0]
  switch(base) {
    case 'host':
      return 'success'
    case 'service':
      return 'warning'
    case 'kubernetes':
      return 'info'
    case 'syslog':
      return 'danger'
    default:
      return 'secondary'
  }
}

function renderAutoComplete(option: SelectOption): VNodeChild {
  return h(NTag,
    {size: 'small', type: selectColor(option.label)},
    {default: () => option.label},
  )
}

</script>

<template>
  <n-card :title="title" size="small">
    <n-grid :cols="1" :y-gap="8">
      <n-gi v-for="[key, val] in dataValue" :key="key">
        <n-input-group>
          <n-input-group-label :size="size">{{ key }}</n-input-group-label>
          <s-multi-value :size="size" :value="val" @update:value="v => update(key, v)" />
          <n-button
            secondary
            type="error"
            :size="size"
            @click="remove(key)"
          >
            <n-icon :component="Trash" />
          </n-button>
        </n-input-group>
      </n-gi>
      <n-gi v-if="dataValue.size == 0">
        <n-empty />
      </n-gi>
    </n-grid>
    <template #action>
      <n-input-group>
        <n-auto-complete
          v-if="autoComplete"
          ref="inputKey"
          v-model:value="newItemKey"
          :options="acOptions"
          :render-label="acRender"
          :size="size"
          placeholder="Key"
        />
        <n-input
          v-else
          ref="inputKey"
          v-model:value="newItemKey"
          :size="size"
          placeholder="Key"
          @keyup.enter="append()"
        />
        <s-multi-value
          v-model:value="newItemValue"
          :size="size"
          @keyup.enter="append()"
        />
        <n-button
          ghost
          type="success"
          :size="size"
          @click="append()"
        >
          <n-icon :component="Add" />
        </n-button>
      </n-input-group>
    </template>
  </n-card>
</template>
