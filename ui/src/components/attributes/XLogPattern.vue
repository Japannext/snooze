<script setup lang="ts">
import { defineProps, h, Component } from 'vue'
import { NDropdown, NButton, NTag, NIcon, NPopover } from 'naive-ui'
import { TypePattern, BookOpen} from '@/icons'

const props = defineProps<{
  name: string,
}>()

const options = [
  {label: 'Documentation', key: 'documentation', icon: renderIcon(BookOpen)},
]

function renderIcon(icon: Component) {
  return () => {
    return h(NIcon, null, {
      default: () => h(icon)
    })
  }
}

function select(key: string) {
}
</script>

<template>
  <n-dropdown
    trigger="click"
    :options="options"
    @select="select"
  >
    <n-button size="tiny">
      <template #icon>
        <n-icon :component="TypePattern" />
      </template>
      {{ props.name }}
    </n-button>
  </n-dropdown>
</template>
