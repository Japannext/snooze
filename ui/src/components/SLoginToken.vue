<script setup lang="ts">
import { ref, defineEmits } from 'vue'
import { NInput, NInputGroup, NButton } from 'naive-ui'

import { LoginService } from '@/api'
import type { TokenCredentials } from '@/api'

const emit = defineEmits<{
  (e: "success", value: string): void,
}>()

const token = ref<string>(null)

function login() {
  const credentials: TokenCredentials = {token: token.value}
  LoginService.authToken(credentials)
  .then((resp) => {
    emit("success", token.value)
  })
  .catch((error) => {
    console.log(JSON.stringify(error))
  })
}

</script>

<template>

  <n-input-group>
    <n-input
      type="password"
      show-password-on="click"
      v-model:value="token"
      round
    />
    <n-button round @click="login">Login</n-button>
  </n-input-group>

</template>
