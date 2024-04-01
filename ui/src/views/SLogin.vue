<script setup lang="ts">

import { NCard, NTabs, NTabPane, NButton } from 'naive-ui'

import { SLoginOidc, SLoginLdap, SLoginToken } from '@/components'
import { OpenAPI } from '@/api'
import { router } from '@/router'
import { xSnoozeToken } from '@/stores'

function authenticated() {
  if ((OpenAPI.HEADERS) && (OpenAPI.HEADERS['X-Snooze-Token'])) {
    return true
  } else {
    return false
  }
}

function setToken(token: string) {
  xSnoozeToken.value = token
  router.back()
}

const defaultLogin = "token"

const loginMethods = {
  token: {title: "Token", kind: "token"},
}

</script>

<template>
  <n-card
    title="Login"
    header-style="text-align: center; padding-top: 5px; padding-bottom: 5px;"
  >
    <n-tabs :default-value="defaultLogin">
      <n-tab-pane
        v-for="(method, name) in loginMethods"
        :key="name"
        :tab="method.title"
        :name="name"
      >
        <s-login-oidc v-if="method.kind == 'oidc'" @success="setToken" />
        <s-login-ldap v-else-if="method.kind == 'ldap'" @success="setToken" />
        <s-login-token v-else @success="setToken" />
      </n-tab-pane>
    </n-tabs>
    <n-button v-if="authenticated()" @click="router.push({name: 'alerts'})">
      Enter
    </n-button>
  </n-card>
</template>
