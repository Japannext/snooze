<script setup lang="ts">

import { withDefaults, computed, defineEmits } from 'vue'
import SConditionChild from '@/components/inputs/SConditionChild.vue'

interface Props {
  value: object,
  size: string,
}

const props = withDefaults(defineProps<Props>(), {
  value: () => {},
  size: "small",
})

const emit = defineEmits<{
  (e: "update:value", value: object): void,
}>()

const dataValue = computed<object>({
  get() {
    console.log(`root dataValue.get(): ${JSON.stringify(props.value)}`)
    return props.value
  },
  set(v) { emit("update:value", v)},
})

</script>

<template>
  <s-condition-child
    v-model:value="dataValue"
    :root="true"
    :size="size"
    :index="0"
  />
</template>
