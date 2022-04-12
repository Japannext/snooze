<template>
  <div>
    <v-swatches
      v-model="dataValue"
      :swatches="swatches"
      show-fallback
      fallback-input-type="color"
      popover-x="left"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { PALETTE } from '@/utils/colors2'

import VSwatches from 'vue3-swatches'
import Base from './Base.vue'

const PALETTE_8 = Array.from(PALETTE.values()).slice(7)

export default defineComponent({
  components: {
    VSwatches,
  },
  extends: Base,
  props: {
    modelValue: {type: String, required: true}
  },
  emits: ['update:modelValue'],
  data() {
    return {
      dataValue: this.modelValue,
      swatches: PALETTE_8,
    }
  },
  watch: {
    dataValue: {
      handler() {
        this.$emit('update:modelValue', this.dataValue)
      },
      immediate: true
    },
  },
})
</script>
