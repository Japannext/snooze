<template>
  <div>
    <CForm class="row g-0" @submit.prevent>
      <CCol xs="auto">
        <CInputGroup>
          <CFormInput
            v-model="dataValue"
            :disabled="disabled"
            aria-describedby="feedback"
            :required="required"
            :invalid="required && !checkField"
            :valid="required && checkField"
            type="number"
            min="-1"
          />
          <CTooltip content="Reset">
            <template #toggler="{ on }">
              <CButton
                :disabled="disabled"
                color="info"
                v-on="on"
                @click="reset"
                @click.stop.prevent
              >
                <i class="la la-redo-alt la-lg"></i>
              </CButton>
            </template>
          </CTooltip>
          <CInputGroupText>{{ converted }}</CInputGroupText>
        </CInputGroup>
      </CCol>
    </CForm>
    <CFormFeedback invalid>Field is required</CFormFeedback>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { prettyDuration } from '@/utils/functions'

import Base from './Base.vue'

const EMPTY_VALUES = [undefined, null, '', [], {}]

// @group Forms
// Class for inputing a duration
export default defineComponent({
  name: 'Duration',
  extends: Base,
  props: {
    modelValue: {type: [String, Number], required: true},
    options: {type: Object, default: () => new Object()},
    disabled: {type: Boolean, default: false},
    required: {type: Boolean, default: false},
    defaultValue: {type: Number, default: 86400},
  },
  emits: ['update:modelValue'],
  data() {
    let value = this.modelValue
    if (EMPTY_VALUES.includes(this.modelValue)) {
      value = this.defaultValue
    }
    return {
      dataValue: value.toString(),
      opts: this.options || {},
    }
  },
  computed: {
    checkField() {
      return this.dataValue != ''
    },
    converted() {
      var dataValue = parseInt(this.dataValue) || 0
      if (dataValue < 0) {
        return this.opts.negative_label || ''
      } else if (dataValue == 0) {
        return this.opts.zero_label || this.opts.negative_label || ''
      } else if (this.opts.custom_label != undefined) {
        return (this.opts.custom_label_prefix || '') + dataValue + (this.opts.custom_label || '')
      } else {
        return (this.opts.custom_label_prefix || '') + this.prettyDuration(dataValue)
      }
    }
  },
  watch: {
    dataValue: {
      handler() {
        this.$emit('update:modelValue', parseInt(this.dataValue) || 0)
      },
      immediate: true
    },
  },
  methods: {
    prettyDuration,
    reset() {
      this.dataValue = this.defaultValue.toString()
    },
  },
})
</script>
