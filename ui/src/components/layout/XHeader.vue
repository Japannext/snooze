<script setup lang="ts">
import { computed } from 'vue'
// import { useLocale } from '@/stores'
import { menuOptions } from '@/router'
import { useRoute } from 'vue-router'

import { NSpace, NMenu, NGrid, NGi} from 'naive-ui'
// import STray from '@/components/STray.vue'

// const locale = useLocale()
/*
const options = [
  {label: "English", value: "en"},
  {label: "日本語", value: "ja"},
]
*/

const route = useRoute()
const selectedMenu = computed(() => {
  if (route.name) {
    return route.name
  }
  return ""
})

/*
const localeData = computed({
  get() { return $i18n.locale },
  set(value) {
    locale.value = value
  },
})
*/

</script>

<template>
  <n-grid layout-shift-disabled>
    <n-gi :span="3">
      <img src="public/img/logo.png" :height="36" style="margin: 2px">
    </n-gi>
    <n-gi :span="18">
      <n-space justify="center">
        <n-menu
          v-model:value="selectedMenu"
          mode="horizontal"
          responsive
          :options="menuOptions"
        />
      </n-space>
    </n-gi>
    <n-gi :span="4">
      <!-- <s-tray /> -->
    </n-gi>
  </n-grid>
  <!-- </n-space> -->
</template>
