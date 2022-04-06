<template>
  <span style="white-space:pre">
    {{ prettyDate(date, showSeconds) }}
  </span>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { prettyDate } from '@/utils/functions'

export default defineComponent({
  props: {
    date: {type: String, required: true},
    showSeconds: {type: Boolean, default: false},
  },
  methods: {
    prettyDate,
  },
})
</script>
