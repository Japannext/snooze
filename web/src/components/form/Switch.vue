<template>
  <div>
    <CForm class="m-0">
      <CFormSwitch
        v-model="dataValue"
        class="pointer"
        :size="size"
        :label="label"
      />
    </CForm>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Base from './Base.vue'

export default defineComponent({
  name: 'Switch',
  extends: Base,
  props: {
    modelValue: {
      type: Boolean,
      default(props: {defaultValue: boolean}) { props.defaultValue },
    },
    defaultValue: {type: Boolean, default: false},
    size: {type: String, default: 'xl'},
    label: {type: String, required: true},
  },
  emits: ['update:modelValue'],
  data() {
    return {
      dataValue: this.modelValue,
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
