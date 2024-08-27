<script setup lang="ts">
import { defineProps, defineExpose } from 'vue'
import { NDropdown } from 'naive-ui'

const props = defineProps<{
  x: number,
  y: number,
  show: boolean
}>()

defineExpose([show])

function show() {

}

</script>

<template>
  <n-dropdown
    placement="bottom-start"
    trigger="manual"
    :x="x"
    :y="y"
    :options="options"
    :show="show"
  />
</template>
