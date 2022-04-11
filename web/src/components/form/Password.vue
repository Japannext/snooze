<template>
  <div>
    <CFormInput
      ref="pwd"
      v-model="dataValue"
      :disabled="disabled"
      type="password"
      :invalid="checkFieldInvalid"
      :valid="checkFieldValid"
    />
    <div class="pt-1">
      <CFormInput
        ref="pwd_confirm"
        v-model="dataValue_repeat"
        :disabled="disabled"
        type="password"
        :invalid="checkFieldInvalid"
        :valid="checkFieldValid"
      />
    </div>
    <CFormFeedback invalid>
      <label v-if="dataValue.length == 0 && datavalue_repeat.length == 0">Fields are required</label>
      <label v-else>Passwords are not identical</label>
    </CFormFeedback>
  </div>
</template>

<script>
// @group Forms
// Class for inputing a string
import Base from './Base.vue'

export default {
  extends: Base,
  props: {
    'modelValue': {type: String, default: () => ''},
    'options': {type: Array, default: () => []},
    'disabled': {type: Boolean, default: () => false},
    'required': {type: Boolean, default: () => false},
  },
  emits: ['update:modelValue'],
  data() {
    return {
      dataValue: this.modelValue || '',
      dataValue_repeat: '',
    }
  },
  computed: {
    checkFieldValid() {
      return this.checkField(true)
    },
    checkFieldInvalid() {
      return !this.checkField(false)
    },
  },
  watch: {
    dataValue: {
      handler: function () {
        this.$emit('update:modelValue', this.dataValue)
      },
      immediate: true
    },
  },
  methods: {
    checkField(is_valid) {
      if (this.dataValue.length == 0 && this.datavalue_repeat.length == 0) {
        return !is_valid && !this.required
      } else {
        return this.dataValue == this.datavalue_repeat
      }
    }
  },
}

</script>
