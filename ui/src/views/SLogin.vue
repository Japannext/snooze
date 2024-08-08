<script setup lang="ts">
import { ref, defineEmits, onMounted, VNodeChild } from 'vue'
import { NSpace, NLayout, NLayoutContent, NButton, NButtonGroup } from 'naive-ui'
import { Openid } from '@vicons/fa'
import { Password16Regular } from '@vicons/fluent'

import axios from 'axios'
import type { AxiosResponse } from 'axios'

import login from '@/api/login'
import type { LoginBackend, LoginBackendKind, LoginBackendList } from '@/api/login'

const backends = ref<LoginBackend[]>([])

onMounted(() => {
  axios.get<LoginBackendList>('/login-backends')
    .then((resp: AxiosResponse<LoginBackendList>) => {
      backends.value = resp.data.items
    })
    .catch((error) => {
      console.log(error)
    })
})

</script>

<template>
  <n-layout>
    <n-layout-content id="container">
      <n-space id="element" vertical>
        <n-button-group>
          <template v-for="(backend, index) in backends" :key="index">
            <n-button v-if="backend.kind == LoginBackendKind.Oidc">
              <template #icon><n-icon :component="Openid" /></template>
              {{ backend.name }}
            </n-button>
            <n-button v-if="backend.kind == LoginBackendKind.Static">
              <template #icon><n-icon :component="Password16Regular" /></template>
              Admin bypass (static)
            </n-button>
          </template>
        </n-button-group>
      </n-space>
    </n-layout-content>
  </n-layout>
</template>
