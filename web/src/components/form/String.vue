<template>
  <div>
    <CFormInput
      v-model="dataValue"
      :disabled="disabled"
      aria-describedby="feedback"
      :invalid="isInvalid"
      :placeholder="placeholder"
    />
    <CFormFeedback invalid>
      Field is required
    </CFormFeedback>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import Base from './Base.vue'

// @group Forms
// Class for inputing a string
export default defineComponent({
  name: 'String',
  extends: Base,
  props: {
    modelValue: {
      type: String as PropType<string|null>,
      default(props: {defaultValue: string|null}) { props.defaultValue },
    },
    defaultValue: {type: String as PropType<string|null>, default: null},
    disabled: {type: Boolean, default: () => false},
    required: {type: Boolean, default: () => false},
    placeholder: {type: String, default: () => ''},
  },
  emits: ['update:modelValue'],
  data() {
    return {
      dataValue: this.modelValue,
    }
  },
  computed: {
    isInvalid () {
      return (this.required && this.dataValue == null)
    },
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
