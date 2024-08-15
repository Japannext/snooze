<script setup lang="ts">
import {
  NLayout, NLayoutContent,
  NLayoutHeader, NLayoutFooter, NLayoutSider,
} from 'naive-ui'

import { useSider } from '@/stores'
import SHeader from '@/components/SHeader.vue'
import SMenu from '@/components/SMenu.vue'
import SFooter from '@/components/SFooter.vue'

const { collapsed } = useSider()

</script>

<template>
  Snooze
  <n-layout position="absolute">
    <n-layout-header
      bordered
      position="absolute"
      style="height: 48px; padding: 12px"
    >
      <s-header />
    </n-layout-header>
    <n-layout
      has-sider
      position="absolute"
      style="top: 48px; bottom: 48px"
    >
      <n-layout-sider
        v-model:collapsed="collapsed"
        bordered
        collapse-mode="width"
        show-trigger
        :width="240"
      >
        <s-menu />
      </n-layout-sider>
      <n-layout-content content-style="padding: 24px;">
        <router-view v-slot="{ Component }">
          <component :is="Component" />
        </router-view>
      </n-layout-content>
    </n-layout>
    <n-layout-footer
      bordered
      position="absolute"
      style="height: 48px; padding: 12px"
    >
      <s-footer />
    </n-layout-footer>
  </n-layout>
</template>

<style scoped>

</style>
