<template>
  <div>
    <CForm inline>
      <CInputGroup class="pe-1">
        <CInputGroupText>Red</CInputGroupText>
        <CFormSelect v-model="dataValue.red" :value="dataValue.red">
          <option v-for="opts in LIGHT_OPTIONS" :key="opts" :value="opts">{{ opts }}</option>
        </CFormSelect>
      </CInputGroup>
      <CInputGroup class="pe-1">
        <CInputGroupText>Yellow</CInputGroupText>
        <CFormSelect v-model="dataValue.yellow" :value="dataValue.yellow">
          <option v-for="opts in LIGHT_OPTIONS" :key="opts" :value="opts">{{ opts }}</option>
        </CFormSelect>
      </CInputGroup>
      <CInputGroup class="pe-1">
        <CInputGroupText>Green</CInputGroupText>
        <CFormSelect v-model="dataValue.green" :value="dataValue.green">
          <option v-for="opts in LIGHT_OPTIONS" :key="opts" :value="opts">{{ opts }}</option>
        </CFormSelect>
      </CInputGroup>
      <CInputGroup class="pe-1">
        <CInputGroupText>Blue</CInputGroupText>
        <CFormSelect v-model="dataValue.blue" :value="dataValue.blue">
          <option v-for="opts in LIGHT_OPTIONS" :key="opts" :value="opts">{{ opts }}</option>
        </CFormSelect>
      </CInputGroup>
      <CInputGroup>
        <CInputGroupText>White</CInputGroupText>
        <CFormSelect v-model="dataValue.white" :value="dataValue.white">
          <option v-for="opts in LIGHT_OPTIONS" :key="opts" :value="opts">{{ opts }}</option>
        </CFormSelect>
      </CInputGroup>
    </CForm>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Base from './Base.vue'

const DEFAULT_PATLITE = {red: 'off', yellow: 'off', green: 'off', blue: 'off', white: 'off'}
const LIGHT_OPTIONS = ['off', 'on', 'blink1', 'blink2']

export default defineComponent({
  name: 'Patlite',
  extends: Base,
  props: {
    modelValue: {type: Object, default: () => DEFAULT_PATLITE },
  },
  emits: ['update:modelValue'],
  data() {
    return {
      LIGHT_OPTIONS: LIGHT_OPTIONS,
      dataValue: this.modelValue,
    }
  },
  watch: {
    dataValue: {
      handler() {
        this.$emit('update:modelValue', this.dataValue)
      },
      deep: true,
      immediate: true
    },
  },
})
</script>
