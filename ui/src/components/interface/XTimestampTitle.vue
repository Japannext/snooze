<script setup lang="ts">
import { defineProps, withDefaults } from 'vue'
import { NSwitch, NSpace } from 'naive-ui'
import { useTimeRelative } from '@/stores'

const { timeRelative } = useTimeRelative()

interface Props {
  title?: string
}

withDefaults(defineProps<Props>(), {
  title: "Timestamp",
})

</script>

<template>
  <n-space>
    <span>{{ title }}</span>
    <n-switch v-model:value="timeRelative" size="small" :round="false" />
  </n-space>
</template>
