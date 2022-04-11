<template>
  <div>
    <CForm @submit.prevent>
      <CRow v-for="(argument, index) in datavalue" :key="index" class="g-0">
        <CCol xs="auto">
          <CInputGroup class="pb-1">
            <CFormInput v-model="argument[0]" :placeholder="placeholder[0]" class="col-form-label" />
            <CFormInput v-model="argument[1]" :placeholder="placeholder[1]" class="col-form-label" />
            <CTooltip content="Remove">
              <template #toggler="{ on }">
                <CButton color="danger" v-on="on" @click="remove(index)" @click.stop.prevent>
                  <i class="la la-trash la-lg"></i>
                </CButton>
              </template>
            </CTooltip>
          </CInputGroup>
        </CCol>
      </CRow>
      <CRow class="g-0">
        <CCol xs="auto">
          <CTooltip content="Add">
            <template #toggler="{ on }">
              <CButton color="secondary" v-on="on" @click="append" @click.stop.prevent>
                <i class="la la-plus la-lg"></i>
              </CButton>
            </template>
          </CTooltip>
        </CCol>
      </CRow>
    </CForm>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import Base from './Base.vue'

export default defineComponent({
  name: 'Arguments',
  extends: Base,
  props: {
    modelValue: {type: Array, default: () => []},
    placeholder: {type: Array, default: () => ['--key', 'value']},
  },
  emits: ['update:modelValue'],
  data () {
    return {
      datavalue: this.modelValue,
    }
  },
  watch: {
    datavalue: {
      handler() {
        this.$emit('update:modelValue', this.datavalue)
      },
      immediate: true,
    },
  },
  methods: {
    append() {
      this.datavalue.push(['', ''])
    },
    remove(index: number) {
      this.datavalue.splice(index, 1)
    }
  },
})
</script>
