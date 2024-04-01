<script setup lang="ts">
import {
  NSpace, NP, NButton, NIcon,
} from 'naive-ui'
import { LogoGithub, Book } from '@vicons/ionicons5'
import { version, homepage, repository } from '../../package'
</script>

<template>
  <n-space justify="space-between" align="center">
    <!-- Left -->
    <span />

    <!-- Center -->
    <n-p>Snooze {{ version }}</n-p>

    <!-- Right -->
    <n-space align="center">
      <a :href="repository.url">
        <n-button text style="font-size: 24px">
          <n-icon :component="LogoGithub" />
        </n-button>
      </a>
      <a :href="homepage">
        <n-button text style="font-size: 24px">
          <n-icon :component="Book" />
        </n-button>
      </a>
    </n-space>
  </n-space>
</template>
