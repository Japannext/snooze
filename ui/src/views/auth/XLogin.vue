<script setup lang="ts">
import { ref, defineEmits, onMounted, VNodeChild } from 'vue'
import { NIcon, NCard, NSpace, NLayout, NLayoutContent, NButton, NButtonGroup } from 'naive-ui'
import { Openid, SignInAlt } from '@vicons/fa'
import { Password16Regular } from '@vicons/fluent'
import { router } from '@/router'

import { getAuthMethods, type AuthMethod, oidcs } from '@/api'

const methods = ref<AuthMethod[]>([])

onMounted(() => {
  getAuthMethods()
    .then((items) => {
      methods.value = items
      oidcs.init(items)
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

function oidcLogin(method: AuthMethod) {
  var backend = oidcs.get(method.name)
  if (backend === undefined) {
    throw new Error(`could not find oidc backend ${method.name}`)
  }
  backend.login()
  .then(() => {
    console.log("logged in successfully")
  })
}

</script>

<template>
  <n-card size="huge" title="Login">
    {{ methods }}
    <n-space id="element" vertical>
      <template v-for="method in methods" :key="method.name">
        <template v-if="method.kind == 'oidc'">
          <n-button :color="method.color" @click="oidcLogin(method)">
            <template v-if="method.icon" #icon><n-icon :component="getIcon(method.icon)" /></template>
            {{ method.displayName }}
          </n-button>
        </template>
        <template v-else>
          <n-button :color="method.color" @click="postLogin(method.name)">
            <template v-if="method.icon" #icon><n-icon :component="getIcon(method.icon)" /></template>
            {{ method.displayName }}
          </n-button>
        </template>
      </template>
    </n-space>
  </n-card>
</template>
