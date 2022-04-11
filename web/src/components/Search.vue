<template>
  <CInputGroup>
    <CFormInput
      v-model="datavalue"
      placeholder="Search"
      type="search"
      class="border-bottom-0"
      style="border-bottom-left-radius: 0"
    />
    <CButton block color="primary" type="submit" @click="search">
      <i class="la la-search la-lg"></i>
    </CButton>
    <CButton
      block
      color="secondary"
      type="reset"
      style="border-bottom-right-radius: 0"
      @click="clear"
    >
      Clear
    </CButton>
    <slot name="search_buttons"></slot>
  </CInputGroup>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    modelValue: {type: String, default: ""},
  },
  emits: ['update:modelValue', 'clear', 'search'],
  data () {
    return {
      datavalue: this.modelValue,
    }
  },
  watch: {
    datavalue () {
      this.$emit('update:modelValue', this.datavalue)
    }
  },
  methods: {
    search() {
      this.$emit('search', this.datavalue)
    },
    clear() {
      this.datavalue = ""
      // Function to call when the clear button is pushed
      this.$emit('clear')
    },
  },
})
</script>
