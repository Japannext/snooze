<script setup lang="ts">
import { defineProps, withDefaults, computed } from 'vue'

import { NTag } from 'naive-ui'

interface Props {
  name: string,
  color?: string,
  closable?: boolean,
}

const props = withDefaults(defineProps<Props>(), {
  color: undefined,
  closable: false,
})

/*
const realTag = computed(() => {
  var cachedTag = tagMap.get(props.tag)
  if (cachedTag) {
    return cachedTag
  }
  return {name: props.tag, color: "", description: ""}
})
*/

const colorStyle = computed(() => {
  if (props.color) {
    return {color: props.color, textColor: '#FFFFFF', borderColor: '#FFFFFF'}
  } else {
    return {color: '#B7B7B7FF', textColor: '#FFFFFF'}
  }
})
</script>

<template>
  <n-tag
    size="small"
    round
    :color="colorStyle"
    :bordered="false"
    :closable="closable"
  >
    {{ name }}
  </n-tag>
</template>

<style scoped>
.n-base-close {
  color: '#FFFFFF';
}
</style>
