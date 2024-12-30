<script setup lang="ts">

import { NMenu } from 'naive-ui'

import { menuOptions } from '@/router'

</script>

<template>
  <n-menu
    inverted
    :options="menuOptions"
  />
</template>
