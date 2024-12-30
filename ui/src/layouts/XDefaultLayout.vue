<script setup lang="ts">

// Components
import { NLayout, NLayoutContent, NLayoutHeader, NLayoutFooter } from 'naive-ui'
import { XHeader, XFooter } from '@/components'

</script>

<template>
  <n-layout position="absolute">
    <n-layout-header
      bordered
      position="absolute"
      style="height: 48px; padding: 3px"
    >
      <x-header />
    </n-layout-header>
    <n-layout
      position="absolute"
      style="top: 48px; bottom: 48px"
    >
      <n-layout-content content-style="padding: 10px;">
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
      <x-footer />
    </n-layout-footer>
  </n-layout>
</template>
