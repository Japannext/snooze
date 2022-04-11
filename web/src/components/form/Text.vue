<template>
  <div>
    <CFormTextarea
      v-model="dataValue"
      name="string"
      :disabled="disabled"
      aria-describedby="feedback"
      :invalid="isInvalid"
      :placeholder="placeholder"
    />
    <CFormFeedback id="feedback" invalid :state="checkField">
      Field is required
    </CFormFeedback>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import Base from './Base.vue'

export default defineComponent({
  name: 'Text',
  extends: Base,
  props: {
    modelValue: {
      type: String as PropType<string|null>,
      default(props: {defaultValue: string|null}) { props.defaultValue },
    },
    defaultValue: {type: String as PropType<string|null>, default: null},
    disabled: {type: Boolean, default: false},
    required: {type: Boolean, default: false},
    placeholder: {type: String, default: ''},
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
      immediate: true
    },
  },
})
</script>
