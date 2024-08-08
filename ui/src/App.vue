<script setup lang="ts">
import '@/utils/token'
import { NMessageProvider } from 'naive-ui'
</script>

<template>
  <n-message-provider>
    <router-view />
  </n-message-provider>
</template>

<style scoped>
</style>
