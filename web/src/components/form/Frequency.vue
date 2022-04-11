<template>
  <div>
    <CRow>
      <CCol>
        <Duration
          v-model="dataValue.delay"
          :default-value="defaultValue.delay"
          :options="OPTIONS.delay"
        />
      </CCol>
      <CCol>
        <Duration
          v-model="dataValue.every"
          :default-value="defaultValue.every"
          :options="OPTIONS.every"
        />
      </CCol>
      <CCol>
        <Duration
          v-model="dataValue.total"
          :default-value="defaultValue.total"
          :options="OPTIONS.total"
        />
      </CCol>
    </CRow>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import Base from './Base.vue'
import Duration from '@/components/form/Duration.vue'

const DEFAULT_VALUE = {every: 0, total: 1, delay: 0}
const EMPTY_VALUES = [undefined, null, '', [], {}]
const OPTIONS = {
  delay: {custom_label_prefix: 'After ', negative_label: 'Immediately'},
  every: {custom_label_prefix: 'Send every ', negative_label: 'Send'},
  total: {custom_label: ' time(s) total', negative_label: 'Forever', zero_label: 'Nothing'},
}

export default defineComponent({
  name: 'Frequency',
  components: {
    Duration,
  },
  extends: Base,
  props: {
    modelValue: {type: Object, default: () => new Object()},
    defaultValue: {type: Object, default: () => {DEFAULT_VALUE}},
  },
  emits: ['update:modelValue'],
  data() {
    let value = this.modelValue
    if (EMPTY_VALUES.includes(this.modelValue)) {
      value = this.defaultValue
    }
    return {
      OPTIONS: OPTIONS,
      dataValue: value,
    }
  },
  watch: {
    dataValue: {
      handler() {
        this.$emit('update:modelValue', this.dataValue)
      },
      immediate: true,
      deep: true,
    },
  },
})
</script>
