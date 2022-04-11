<template>
  <div>
    <CForm class="m-0">
      <CButtonGroup role="group">
        <CFormCheck
          v-for="(opts, i) in options"
          :id="id + i"
          :key="opts.value"
          type="radio"
          autocomplete="off"
          :name="id"
          :button="{color: 'primary', variant: 'outline'}"
          :checked="opts.value != undefined ? (opts.value == dataValue) : (opts == dataValue)"
          :value="opts.value != undefined ? opts.value : opts"
          :label="opts.text != undefined ? opts.text : opts"
          @click="opts.value != undefined ? (dataValue = opts.value) : (dataValue = opts)"
        />
      </CButtonGroup>
    </CForm>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { SelectorOptions } from '@/utils/types'

import Base from './Base.vue'

// Create a selector form
export default defineComponent({
  name: 'Radio',
  extends: Base,
  props: {
    id: {type: String, required: true},
    modelValue: {type: String, default: null},
    options: {type: Array as PropType<SelectorOptions>, required: true},
  },
  emits: ['update:modelValue'],
  data() {
    let value = this.modelValue
    if (value == null) {
      value = this.options[0].value as string
    }
    return {
      dataValue: value,
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
