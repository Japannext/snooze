<template>
  <span v-if="data === undefined">Empty</span>
  <span v-else>
    <CBadge v-for="(field, index) in data" :key="field" :class="index != data.length - 1 ? 'me-1' : ''" :color="get_color(colorize?field:'')">{{ field }}</CBadge>
  </span>
</template>

<script>

import { get_color } from '@/utils/colors'

export default {
  name: 'Field',
  props: {
    // Array representing the condition
    data: {type: Array},
    // Assign random colors
    colorize: {type: Boolean},
  },
  data() {
    return {
      get_color: get_color,
    }
  },
}
</script>
