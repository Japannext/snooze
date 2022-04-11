<template>
  <div>
    <CFormInput
      v-model="dataValue"
      type="number"
      :disabled="disabled"
      aria-describedby="feedback"
      :required="required"
      :invalid="required && !checkField"
      :valid="required && checkField"
      :min="opts.min"
      :max="opts.max"
    />
    <CFormFeedback invalid>
      Field is required
    </CFormFeedback>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import Base from './Base.vue'

// @group Forms
// Class for inputing a number
export default defineComponent({
  extends: Base,
  props: {
    modelValue: {type: [String, Number], default: () => 0},
    options: {type: Array, default: () => []},
    disabled: {type: Boolean, default: () => false},
    required: {type: Boolean, default: () => false},
    defaultValue: {type: Number, default: () => 0},
  },
  emits: ['update:modelValue'],
  data() {
    return {
      dataValue: ([undefined, 0, [], {}].includes(this.modelValue) ? (this.default_value == undefined ? 0 : this.default_value) : this.modelValue).toString(),
      opts: this.options || {},
    }
  },
  computed: {
    checkField () {
      return this.dataValue != ''
    },
  },
  watch: {
    dataValue: {
      handler: function () {
        this.$emit('update:modelValue', parseInt(this.dataValue) || 0)
      },
      immediate: true
    },
  },
})
</script>
