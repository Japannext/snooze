<script setup lang="ts">
import { NSwitch, NSpace } from 'naive-ui'
import { useTimeRelative } from '@/stores'

const { timeRelative } = useTimeRelative()

</script>

<template>
  <n-space>
    <span>Timestamp</span>
    <n-switch v-model:value="timeRelative" size="small" :round="false">
    </n-switch>
  </n-space>
</template>
