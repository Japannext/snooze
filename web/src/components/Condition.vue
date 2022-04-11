<template>
  <span v-safe-html="htmlValue" />
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { ConditionObject } from '@/utils/condition2'

export default defineComponent({
  name: 'Condition',
  props: {
    // Array representing the condition
    data: {type: Array, required: true},
  },
  computed: {
    htmlValue () {
      return ConditionObject.fromArray(this.data).toHTML()
    },
  },
})
</script>
