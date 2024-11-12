<script setup lang="ts">
import { defineModel } from 'vue'
import { useMouse } from '@vueuse/core'
// Components
import { NDropdown } from 'naive-ui'

const {x, y } = useMouse()

const show = defineModel<boolean>('show', {default: false})

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
