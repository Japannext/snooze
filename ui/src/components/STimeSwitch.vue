<script setup lang="ts">
import type { CSSProperties } from 'vue'
import { NSwitch, NSpace, NIcon } from 'naive-ui'
import { useTimeRelative } from '@/stores'
import { World } from '@/icons'

const { timeRelative } = useTimeRelative()

function railStyle(focused: boolean, checked: boolean): CSSProperties {
  const style: CSSProperties = {}
  if (checked) {
    style.background = ""
  }
  else {
    
  }
  return style
}

</script>

<template>
  <n-switch v-model:value="timeRelative" size="small" :round="false">
    <template #unchecked>ISO</template>
    <template #checked>Relative</template>
  </n-switch>
</template>
