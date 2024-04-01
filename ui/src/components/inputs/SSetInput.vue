<script setup lang="ts">
import { ref, computed, watch, nextTick, defineEmits } from 'vue'
import { useVModel } from '@vueuse/core'
import { NDynamicTags, NAutoComplete, NButton, NIcon } from 'naive-ui'
import { Add } from '@vicons/ionicons5'

const props = defineProps<{
  value: Array<string>,
  options: Array<string>,
  size: string,
}>()

const emit = defineEmits<{
  (e: "update:value", value: Array<string>): void,
}>()

const dataValue = useVModel(props, 'value', emit)

const autoCompleteRef = ref(null)
const inputValue = ref("")

const dataOptions = computed(() => {
  if (inputValue.value == "") {
    return props.options
  }
  return props.options.filter((option) => {
    return option.label.startsWith(inputValue.value)
  })
})

watch(autoCompleteRef, (value) => {
  if (value) {
    nextTick(() => value.focus())
  }
})

</script>
<template>
  <n-dynamic-tags v-model:value="dataValue">
    <template #input="{ submit, deactivate }">
      <n-auto-complete
        ref="autoCompleteRef"
        v-model:value="inputValue"
        :options="dataOptions"
        :size="size"
        :clear-after-select="true"
        @select="submit($event)"
        @blur="deactivate"
      />
    </template>
    <template #trigger="{ activate, disabled }">
      <n-button
        type="primary"
        dashed
        :size="size"
        :disabled="disabled"
        @click="activate()"
      >
        <template #icon><n-icon :component="Add" /></template>
        Add
      </n-button>
    </template>
  </n-dynamic-tags>
</template>
