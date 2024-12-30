<script setup lang="ts">
import { defineProps, ref, onMounted } from 'vue'
import { NTag } from 'naive-ui'
import { type NaiveColor } from '@/api'

const props = defineProps<{
  text: string|undefined,
  number: number,
}>()

const color = ref<NaiveColor>("default")

onMounted(() => {
  if (props.number == 0) { // Not set
    color.value = "default"
  }
  else if (props.number < 5) { // TRACE
    color.value = "default"
  }
  else if (props.number < 9) { // DEBUG
    color.value = "default"
  }
  else if (props.number < 13) { // INFO
    color.value = "info"
  }
  else if (props.number < 17) { // WARN
    color.value = "warning"
  }
  else if (props.number < 21) { // ERROR
    color.value = "error"
  }
  else if (props.number < 25) { // FATAL
    color.value = "error"
  }
  else { // UNKNOWN
    color.value = "default"
  }
})

</script>
<template>
  <n-tag v-if="text" size="small" :type="color">
    {{ text }}
  </n-tag>
</template>
