<template>
  <div>
    <CRow>
      <CCol>
        <VueDatePicker
          v-model="datavalue"
          format="yyyy-MM-dd HH:mm"
          preview-format="yyyy-MM-dd HH:mm"
          placeholder="Select date"
          :input-class-name="datavalue != null ? 'form-control is-valid' : 'form-control is-invalid'"
          :week-start="weekStart"
          :close-on-auto-apply="false"
          text-input
          auto-apply
        />
      </CCol>
    </CRow>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import moment from 'moment'

import VueDatePicker from '@vuepic/vue-datepicker'
import Base from './Base.vue'

export default defineComponent({
  name: 'DateTimeSingle',
  components: {
    VueDatePicker,
  },
  extends: Base,
  props: {
    modelValue: {type: String, default: () => moment().format()},
  },
  emits: ['update:modelValue'],
  data() {
    return {
      dataValue: this.modelValue || moment().format(),
      weekStart: moment().startOf('week').weekday(),
    }
  },
  computed: {
    formattedDate () {
       if (this.dataValue != null) {
         return moment(this.dataValue).format("YYYY-MM-DDTHH:mmZ")
       } else {
         return ''
       }
    }
  },
  watch: {
    dataValue: {
      handler() {
        this.$emit('update:modelValue', this.formattedDate)
      },
      immediate: true,
      deep: true,
    },
  },
  mounted() {
    this.dataValue = this.modelValue || moment().format()
  },
})
</script>

<style lang="scss">
@import '@vuepic/vue-datepicker/src/VueDatePicker/style/main.scss'
</style>
