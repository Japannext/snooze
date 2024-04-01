<script setup lang="ts">
import { withDefaults } from 'vue'
import { useClipboard, usePermission } from '@vueuse/core'
import { NPopover, NButton, NIcon, useMessage } from 'naive-ui'
import { Hashtag } from '@vicons/fa'

const message = useMessage()
const { copy } = useClipboard()
const permissionWrite = usePermission('clipboard-write')

const props = withDefaults(defineProps<{
  hash: string,
}>(), {
  hash: "",
})

function copyHash() {
  copy(props.hash)
  if (permissionWrite.value) {
    message.info(`Copied '${props.hash}' to clipboard`)
  } else {
    message.error(`The browser do not have the permission to copy the hash to clipboard.`)
  }
}

</script>
<template>
  <n-popover trigger="hover" placement="left">
    <template #trigger>
      <n-button circle :size="size" @click="copyHash()">
        <template #icon><n-icon :component="Hashtag" /></template>
      </n-button>
    </template>
    <pre>Hash: {{ hash }}</pre>
  </n-popover>
</template>
