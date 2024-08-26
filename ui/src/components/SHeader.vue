<script setup lang="ts">
import { computed } from 'vue'
import { NSpace, NSwitch, NSelect, NMenu, NFlex, NGrid, NGi, NInputGroup, NButton} from 'naive-ui'
import { useDevMode, useLocale } from '@/stores'
import { menuOptions } from '@/router'
import STray from '@/components/STray.vue'

const { locale } = useLocale()
const { devMode } = useDevMode()
const allowDevMode = Boolean(import.meta.env.MODE == "development")
const options = [
  {label: "English", value: "en"},
  {label: "日本語", value: "ja"},
]

const localeData = computed({
  get() { return $i18n.locale },
  set(value) {
    locale.value = value
  },
})

</script>

<template>
  <!--
  <n-space align="center">
    <img
      name="logo"
      src="img/logo.png"
      :height="24"
    />
    -->
    <!--
    <n-switch v-if="allowDevMode" v-model:value="devMode">
      <template #checked>dev mode</template>
      <template #unchecked>dev mode</template>
    </n-switch>
    <n-select v-model:value="$i18n.locale" :options="options" />
    -->
  <n-grid layout-shift-disabled>
    <n-gi :span="4">
      Snooze
    </n-gi>
    <n-gi :span="16">
      <n-space justify="center">
      <n-menu
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
