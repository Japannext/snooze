<script setup lang="ts">
import { defineProps, withDefaults, computed } from 'vue'

import { NTag } from 'naive-ui'

import type { Tag } from '@/api'

interface Props {
  tag: Tag,
  closable?: boolean,
}

const props = withDefaults(defineProps<Props>(), {
  closable: false,
})

const colorStyle = computed(() => {
  if (props.tag) {
    return {color: props.tag.color, textColor: '#FFFFFF', borderColor: '#FFFFFF'}
  } else {
    return {color: '#B7B7B7FF', textColor: '#FFFFFF'}
  }
})

</script>

<template>
  <n-tag
    v-if="tag"
    size="small"
    round
    :color="colorStyle"
    :bordered="false"
    :closable="closable"
  >
    {{ tag.name }}
  </n-tag>
  <n-tag v-else>...</n-tag>
</template>

<style scoped>
.n-base-close {
  color: '#FFFFFF';
}
</style>
