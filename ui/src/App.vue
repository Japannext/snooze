<script setup lang="ts">
import '@/utils/token'
import { NConfigProvider, NDialogProvider, NNotificationProvider, NMessageProvider, NLoadingBarProvider } from 'naive-ui'
</script>

<template>
  <n-config-provider
    class="demo"
    :locale="lang"
    :theme="theme"
    :theme-overrides="themeOverrides"
  >
    <n-dialog-provider>
      <n-message-provider>
        <n-notification-provider>
          <n-loading-bar-provider>
            <router-view />
          </n-loading-bar-provider>
        </n-notification-provider>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<style scoped>
</style>
