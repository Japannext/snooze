<template>
  <CBadge :style="genColor(data)">
    {{ data }}
  </CBadge>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { genColor } from '@/utils/colors'

export default defineComponent({
  name: 'ColorBadge',
  props: {
    data: {type: String, required: true},
  },
  data() {
    return {
      genColor: genColor,
    }
  },
})
</script>
