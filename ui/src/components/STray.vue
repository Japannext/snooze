<script setup lang="ts">
import { NInputGroup, NButton, NSwitch, NSpace } from 'naive-ui'
import { useTimeRelative } from '@/stores'

const { timeRelative } = useTimeRelative()

</script>

<template>
  <n-space align="center" justify="end" style="padding: 5px;">
    <n-switch v-model:value="timeRelative">
      <template #checked>
        Relative time
      </template>
      <template #unchecked>
        ISO time
      </template>
    </n-switch>
    <n-button>Login</n-button>
  </n-space>
</template>
