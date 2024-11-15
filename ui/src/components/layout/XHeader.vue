<script setup lang="ts">
import { ref, computed } from 'vue'
import { useDevMode, useLocale } from '@/stores'
import { menuOptions } from '@/router'

import { NSpace, NMenu, NFlex, NGrid, NGi, NInputGroup, NButton} from 'naive-ui'
import STray from '@/components/STray.vue'

import type { MenuInst, MenuOption } from 'naive-ui'

const { locale } = useLocale()
const { devMode } = useDevMode()
const allowDevMode = Boolean(import.meta.env.MODE == "development")
const options = [
  {label: "English", value: "en"},
  {label: "日本語", value: "ja"},
]
const selectedMenu = ref<string>("logs")
const menuInstRef = ref<MenuInst | null>(null)

const localeData = computed({
  get() { return $i18n.locale },
  set(value) {
    locale.value = value
  },
})

</script>

<template>
  <n-grid layout-shift-disabled>
    <n-gi :span="2">
      Snooze
    </n-gi>
    <n-gi :span="18">
      <n-space justify="center">
        <n-menu
          ref="nMenu"
          v-model:value="selectedMenu"
          mode="horizontal"
          responsive
          :options="menuOptions"
        />
      </n-space>
    </n-gi>
    <n-gi :span="4">
      <s-tray />
    </n-gi>
  </n-grid>
  <!-- </n-space> -->
</template>
