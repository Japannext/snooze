<script setup lang="ts">
import { defineProps } from 'vue'
import { NTag, NIcon } from 'naive-ui'
import { Mail, LogoGoogle, BuildingLighthouse20Filled, Box24Regular, Browser } from '@/icons'

import type { Destination } from '@/api/types'

const props = defineProps<{
  destination: Destination,
}>()

</script>

<template>
  <div>
    <n-tag v-if="destination.queue == 'mail'">
      <template #icon><n-icon :component="Mail" /></template>
      {{ destination.profile }}
    </n-tag>
    <n-tag v-else-if="destination.queue == 'googlechat'">
      <template #icon><n-icon :component="LogoGoogle" /></template>
      {{ destination.profile }}
    </n-tag>
    <n-tag v-else-if="destination.queue == 'patlite'">
      <template #icon><n-icon :component="BuildingLighthouse20Filled" /></template>
      {{ destination.profile }}
    </n-tag>
    <n-tag v-else-if="destination.queue == 'slack'">
      <template #icon><n-icon :component="Slack" /></template>
      {{ destination.profile }}
    </n-tag>
    <n-tag v-else-if="destination.queue == 'browser'">
      <template #icon><n-icon :component="Browser" /></template>
      {{ destination.profile }}
    </n-tag>
    <n-tag v-else-if="destination.queue == 'dummy'">
      <template #icon><n-icon :component="Box24Regular" /></template>
      {{ destination.profile }}
    </n-tag>
  </div>
</template>
