<script setup lang="ts">
import { ref, onMounted, defineEmits } from 'vue'

import { NInput, NInputGroup, NButton, NIcon } from 'naive-ui'
import { Search } from '@vicons/ionicons5'
import { useRouteQuery } from '@vueuse/router'

const emit = defineEmits<{
  "search": [value: string]
}>()

const search = useRouteQuery("search", "")
const tmp = ref<string>("")

onMounted(() => {
  tmp.value = search.value
})

function run() {
  search.value = tmp.value
  emit("search", tmp.value)
}
</script>

<template>
  <n-input-group>
    <n-input
      v-model:value="tmp"
      placeholder="Search"
      clearable
      @keyup.enter="run()"
    />
    <n-button type="info" @click="run()"><n-icon :component="Search" /></n-button>
  </n-input-group>
</template>
