<script setup lang="ts">

import {
  ref,
  defineEmits,
} from 'vue'

import {
  NInput, NInputGroup,
  NButton, NIcon,
} from 'naive-ui'

import { Search, Close } from '@vicons/ionicons5'

const emit = defineEmits<{
  // Pass the query to parent when the search bar is trigerred
  (e: "search", value: string): void,
}>()

const query = ref<string>("")

// Emit the query to parent
function search() {
  emit("search", query)
}

// Clear the search bar and emit to parent
function clear() {
  query.value = ""
  emit("search", "")
}

</script>

<template>
  <n-input-group>
    <n-input
      v-model:value="query"
      clearable
      placeholder="Search"
      @keyup.enter="search"
    />
    <n-button secondary type="info" @click="search"><n-icon :component="Search" /></n-button>
    <n-button secondary @click="clear"><n-icon :component="Close" /></n-button>
  </n-input-group>
</template>
