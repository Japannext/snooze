<script setup lang="ts">

import {
  NSpace,
  NButton, NIcon,
  NPopover, NEllipsis,
  NDivider,
  NTable, NEmpty,
} from 'naive-ui'

import { List } from '@vicons/ionicons5'

interface Props {
  attributes: object
  size: string
}

const props = defineProps<Props>()

function getVariant(key: string) {
  if (key.startsWith("syslog.")) {
    return "warning"
  } else if (key.startsWith("k8s.")) {
    return "info"
  } else {
    return "default"
  }
}

</script>

<template>
  <n-popover trigger="click">
    <template #trigger>
      <n-button round :size="size">
        <template #icon><n-icon :component="List" /></template>
        {{ Object.keys(attributes).length }}
      </n-button>
    </template>
    <n-table v-if="Object.keys(attributes).length > 0">
      <thead>
        <th>Key</th>
        <th>Value</th>
      </thead>
      <tbody>
        <tr v-for="(value, key) in attributes" :key="key">
          <td>{{ key }}</td>
          <td>{{ value }}</td>
        </tr>
      </tbody>
    </n-table>
    <n-empty small v-else />
  </n-popover>
</template>
