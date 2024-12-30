<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NIcon, NCard, NSpace, NButton } from 'naive-ui'
import { Github, Openid, SignInAlt } from '@/icons'

import { getAuthConfig, type AuthConfig } from '@/api'

const authConfig = ref<AuthConfig>()

onMounted(() => {
  getAuthConfig()
    .then((cfg) => {
      authConfig.value = cfg
    })
})

function getIcon(icon: string) {
  switch (icon.toLowerCase()) {
  case "openid":
    return Openid
  default:
    return SignInAlt
  }
}

function redirect(providerName: string) {
  window.location.href = `/api/auth/${providerName}/login`
}

</script>

<template>
  <n-space vertical style="width: 100">
    <img src="img/logo.png" :width="500" />
    <n-card v-if="authConfig" size="huge" title="Login">
      <n-space vertical>
        <n-button
          v-if="authConfig.oidc"
          :color="authConfig.oidc.color"
          @click="redirect('oidc')"
        >
          <template #icon><n-icon :component="getIcon(authConfig.oidc.icon)" /></template>
          {{ authConfig.oidc.displayName }}
        </n-button>
      </n-space>
    </n-card>
    <n-card v-else>
      Could not get authConfig
    </n-card>
  </n-space>
</template>
