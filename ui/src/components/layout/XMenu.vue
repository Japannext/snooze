<script setup lang="ts">

import { NMenu } from 'naive-ui'

import { menuOptions } from '@/router'

console.log(`menuOptions: ${JSON.stringify(menuOptions)}`)

</script>

<template>
  <n-menu
    :options="menuOptions"
  />
</template>
