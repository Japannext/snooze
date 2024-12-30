<script setup lang="ts">
import { LogoGithub, Book } from '@vicons/ionicons5'
// import { version, homepage, repository } from '../../../package.json'
// import { version, homepage, repository } from '@/../package.json'
// import { version, homepage, repository } from 'pkginfo'

const version = "x.x.x"
const homepage = "http://github.com/japannext/snooze"
const repository = {url: "http://github.com/japannext/snooze"}

// Components
import { NSpace, NP, NButton, NIcon, NButtonGroup, NPopover } from 'naive-ui'
import { XPatlite } from '@/components'
</script>

<template>
  <n-space justify="space-between" align="center">
    <!-- Left -->
    <n-button-group>
      <n-popover trigger="click">
        <template #trigger>
          <n-button>Patlite</n-button>
        </template>
        <x-patlite />
      </n-popover>
    </n-button-group>

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
