<script setup lang="ts">
import { ref, onMounted, defineEmits, defineModel, type Ref } from 'vue'

import { NInput, NInputGroup, NButton, NIcon } from 'naive-ui'
import { Search } from '@vicons/ionicons5'
import { useRouteQuery } from '@vueuse/router'

const route = useRouteQuery<string>("search", "")
const tmp = ref<string>("")

onMounted(() => {
  tmp.value = route.value
})

defineProps<{
  value: string|undefined,
}>()

const emit = defineEmits(['update:value', 'change'])

function run() {
  route.value = tmp.value
  emit('update:value', tmp.value)
  emit('change')
}
</script>

<template>
  <n-input-group>
    <n-input
      v-model:value="tmp"
      placeholder="Search"
      clearable
      @keyup.enter="run"
    />
    <n-button type="info" @click="run"><n-icon :component="Search" /></n-button>
  </n-input-group>
</template>
