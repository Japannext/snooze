<script setup lang="ts">

import { ref, onMounted, } from 'vue'
import {
  NEmpty,
} from 'naive-ui'

interface Props {
  oid: string,
  service,
}

const props = defineProps<Props>()

const item = ref<object>(null)

onMounted(() => {
  props.service.get(props.oid)
  .then((obj) => {
    item.value = obj
  })
})

</script>

<template>
  <n-empty v-if="item == 'null'" />
  <pre v-else>{{ item }}</pre>
</template>
