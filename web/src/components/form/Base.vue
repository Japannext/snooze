<template>
  <div>
    <CRow>
      <CCol col=3 md=2>
        <label
          :id="`title_${metadata.display_name}`"
          v-c-tooltip="{content: this.metadata.description, placement: 'right'}"
        >
          {{ metadata.display_name }}
        </label>
        <label v-if="metadata.required">*</label>
      </CCol>
      <CCol col=9 md=10>
        <component
          :is="component"
          :id="`component_${metadata.display_name}`"
          v-model="dataValue"
          :metadata="metadat"
        />
      </CCol>
    </CRow>
  </div>
</template>

<script lang="ts">
import { defineAsyncComponent, defineComponent, shallowRef } from 'vue'

// @group Forms
// Base class for all form inputs
export default defineComponent({
  props: {
    modelValue: {},
    metadata: {type: Object, default: () => new Object()},
    data: {type: Object},
    required: {type: Boolean, default: false}
  },
  emits: ['update:modelValue'],
  data() {
    return {
      dataValue: (this.modelValue != undefined) ? this.modelValue : (this.metadata ? this.metadata.default : {}),
      component: shallowRef(defineAsyncComponent(() => import(`./${this.metadata.component}.vue`))),
      popover: {content: this.metadata ? (this.metadata.description || '') : '', trigger: ['hover', 'focus'], placement: 'right'}
    }
  },
  watch: {
    dataValue () {
      // Return the value of the input form
      this.$emit('update:modelValue', this.dataValue)
    }
  },
})
</script>
