<script setup lang="ts">
import { NSpace, NTag } from 'naive-ui'
import { Hashtag } from '@vicons/fa'

import type { Group } from '@/api'
import SResource from '@/components/SResource.vue'
import SHash from '@/components/SHash.vue'

const props = defineProps<{
  group: Group,
  size: string,
}>()

function unwrap(keyvalues: object): object {
  const newObject = new Object()
  for (const [key, value] of Object.entries(keyvalues)) {
    const newKey = key.split('.').splice(1).join('.')
    newObject[newKey] = value
  }
  return newObject
}
</script>

<template>
  <n-space align="center">
    <s-hash :hash="group.hash" />
    <s-resource
      :resource="unwrap(group.attributes)"
      :size="size"
    />
  </n-space>
</template>
