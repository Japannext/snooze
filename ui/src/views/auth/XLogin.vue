<script setup lang="ts">
import { ref, defineEmits, onMounted, VNodeChild } from 'vue'
import { NIcon, NCard, NSpace, NLayout, NLayoutContent, NButton, NButtonGroup } from 'naive-ui'
import { router } from '@/router'
import { useRouter } from 'vue-router'
import { Github, Openid, SignInAlt } from '@/icons'

import { getAuthConfig, type AuthConfig } from '@/api'

const router = useRouter()
const authConfig = ref<AuthConfig>(null)

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

function githubLogin() {
  window.location.href = `/api/auth/github/login`
}

function oidcLogin(provider: string) {
  // router.push({path: `/api/auth/${provider}/login`})
  window.location.href = `/api/auth/${provider}/login`
}

</script>

<template>
  <n-card size="huge" title="Login" v-if="authConfig">
    {{ authConfig }}

    <n-space vertical>
      <n-button
        v-if="authConfig.oidc"
        :color="authConfig.oidc.color"
        @click="oidcLogin('openid-connect')"
      >
        <template #icon><n-icon :component="getIcon(authConfig.oidc.icon)" /></template>
        {{ authConfig.oidc.displayName }}
      </n-button>

      <n-button
        v-if="authConfig.github"
        @click="githubLogin()"
      >
        <template #icon><n-icon :component="Github" /></template>
        GitHub
      </n-button>
    </n-space>
  </n-card>
  <n-card v-else>
    Could not get authConfig
  </n-card>

</template>
