<template>
  <div>
    <CFormSelect
      v-model="dataValue"
      :invalid="isInvalid"
    >
      <option disabled :value="null">
        Please select an option
      </option>

      <option v-for="opt in options" :key="opt.value" :value="opt.value">
        {{ opts.text }}
      </option>
    </CFormSelect>

    <CFormFeedback invalid>
      Field is required
    </CFormFeedback>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { SelectorValue, SelectorOptions } from '@/utils/types'

import Base from './Base.vue'

// Create a selector form
export default defineComponent({
  name: 'Selector',
  extends: Base,
  props: {
    modelValue: {
      type: Object as PropType<SelectorValue|null>,
      default(props: {defaultValue: SelectorValue}) { props.defaultValue },
    },
    options: {type: Array as PropType<SelectorOptions>, required: true},
    defaultValue: {type: Object as PropType<SelectorValue>, default: null},
    required: {type: Boolean, default: false},
  },
  emits: ['update:modelValue'],
  data() {
    return {
      dataValue: this.modelValue,
    }
  },
  computed: {
    isInvalid() {
      return (this.required && this.dataValue == null)
    },
  },
  watch: {
    dataValue: {
      handler() {
        this.$emit('update:modelValue', this.dataValue)
      },
      immediate: true,
    },
  },
})
</script>
