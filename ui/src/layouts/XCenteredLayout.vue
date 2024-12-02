<script setup lang="ts">
import {
  NSpace, NLayout, NLayoutContent,
} from 'naive-ui'

</script>

<template>
  <n-layout>
    <n-layout-content id="container">
      <n-space id="element" vertical>
        <router-view v-slot="{ Component }">
          <component :is="Component" id="component" />
        </router-view>
      </n-space>
    </n-layout-content>
  </n-layout>
</template>

<style scoped>
.n-layout-content {
  margin: 0;
  display: flex;
  place-items: center;
  justify-content: center;
  min-width: 320px;
  min-height: 100vh;
}

#element {
  max-width: 1280px;
  min-width: 40vw;
  margin: 0 auto;
  padding: 2rem;
}
</style>
