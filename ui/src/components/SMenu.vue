<script setup lang="ts">

import { NMenu } from 'naive-ui'

import { menuOptions } from '@/router'

console.log(`menuOptions: ${JSON.stringify(menuOptions)}`)

</script>

<template>
  <div>
    <n-menu
      bordered
      :options="menuOptions"
    />
  </div>
</template>
