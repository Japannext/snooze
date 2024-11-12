<script setup lang="ts">
import { ref, defineEmits, onMounted, VNodeChild } from 'vue'
import { NIcon, NCard, NSpace, NLayout, NLayoutContent, NButton, NButtonGroup } from 'naive-ui'
import { Openid, SignInAlt } from '@vicons/fa'
import { Password16Regular } from '@vicons/fluent'
import { xSnoozeToken } from '@/stores'
import { router } from '@/router'

import axios from 'axios'
import type { AxiosResponse } from 'axios'

import { getAuthButtons } from '@/api/login'
import type { AuthButtons } from '@/api/login'

const buttons = ref<AuthButtons[]>([])

onMounted(() => {
  getAuthButtons().then((resp: AxiosResponse<AuthButtons>) => {
    buttons.value = resp.data
  })
})

function getComponent(icon: string) {
  switch (icon) {
    case "openid": {
      return Openid
    }
    default: {
      return SignInAlt
    }
  }
}

function redirect(path: string) {
  router.push({path: path})
}

</script>

<template>
  <n-card size="huge" title="Login">
    {{ buttons }}
    <n-space id="element" vertical>
      <n-button-group>
        <template v-for="(button, index) in buttons" :key="index">
          <a :href="button.path">
            <n-button :color="button.color" @click="redirect(button.path)">
              <template v-if="button.icon" #icon>
                <n-icon :component="getComponent(button.icon)" />
              </template>
              {{ button.displayName }}
            </n-button>
          </a>
        </template>
      </n-button-group>
    </n-space>
  </n-card>
</template>
