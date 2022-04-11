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

import { gen_palette } from '@/utils/colors'

import VSwatches from 'vue3-swatches'

import Base from './Base.vue'

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
      swatches: [],
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
  mounted() {
    this.swatches = gen_palette(8, 8)
  },
})
</script>
