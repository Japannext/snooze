<script setup lang="ts">
import { computed } from 'vue'
import { NSpace, NSwitch, NSelect } from 'naive-ui'
import { useDevMode, useLocale } from '@/stores'

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
  <n-space align="center" inline>
    <img
      name="logo"
      src="img/logo.png"
      :height="24"
    />
    <n-switch v-if="allowDevMode" v-model:value="devMode">
      <template #checked>dev mode</template>
      <template #unchecked>dev mode</template>
    </n-switch>
    <n-select v-model:value="$i18n.locale" :options="options" />
  </n-space>
</template>
