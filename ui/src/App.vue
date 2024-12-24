<script setup lang="ts">
import { NConfigProvider, NDialogProvider, NNotificationProvider, NMessageProvider, NLoadingBarProvider, NThemeEditor } from 'naive-ui'
import '@/utils/token'
</script>

<template>
  <n-config-provider
    class="config"
  >
    <n-dialog-provider>
      <n-message-provider>
        <n-notification-provider>
          <n-loading-bar-provider>
            <n-theme-editor>
              <router-view />
            </n-theme-editor>
          </n-loading-bar-provider>
        </n-notification-provider>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<style scoped>
</style>
