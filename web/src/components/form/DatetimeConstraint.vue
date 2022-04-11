<template>
  <div>
    <VueDatePicker
      v-model="dataValue"
      format="FRONTEND_TIME_FORMAT"
      preview-format="FRONTEND_TIME_FORMAT"
      placeholder="Select date"
      :input-class-name="dataValue != null ? 'form-control is-valid' : 'form-control is-invalid'"
      :week-start="weekStart"
      :close-on-auto-apply="false"
      text-input
      auto-apply
      range
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import moment from 'moment'
import { DatetimeConstraint } from '@/utils/types'

import VueDatePicker from '@vuepic/vue-datepicker';
import Base from './Base.vue'

const FRONTEND_TIME_FORMAT = 'yyyy-MM-dd HH:mm'
const BACKEND_TIME_FORMAT = 'YYYY-MM-DDTHH:mmZ'

function defaultDatetime(): DatetimeConstraint {
  const now = moment()
  const oneHourLater = now.clone().add(1, 'hours')
  return {
    from: now.format(BACKEND_TIME_FORMAT),
    until: oneHourLater.format(BACKEND_TIME_FORMAT),
  }
}

export default defineComponent({
  name: 'DatetimeConstraint',
  components: {
    VueDatePicker,
  },
  extends: Base,
  props: {
    modelValue: {type: Object as PropType<DatetimeConstraint>, default: defaultDatetime},
  },
  emits: {
    'update:modelValue': (datetime: DatetimeConstraint) => {return true},
  },
  data() {
    const now = moment()
    const oneHourLater = now.clone().add(1, 'hours')
    const from = this.modelValue['from'] || now.format(BACKEND_TIME_FORMAT)
    const until = this.modelValue['until'] || oneHourLater.format(BACKEND_TIME_FORMAT)
    return {
      FRONTEND_TIME_FORMAT: FRONTEND_TIME_FORMAT,
      BACKEND_TIME_FORMAT: BACKEND_TIME_FORMAT,
      weekStart: now.startOf('week').weekday(),
      dataValue: [from, until],
    }
  },
  watch: {
    dataValue: {
      handler() {
        const [from, until] = this.dataValue
        const formattedDate = {
          from: moment(from).format(BACKEND_TIME_FORMAT),
          until: moment(until).format(BACKEND_TIME_FORMAT),
        }
        this.$emit('update:modelValue', formattedDate)
      },
      immediate: true,
      deep: true,
    },
  },
})
</script>

<style lang="scss">
@import '@vuepic/vue-datepicker/src/VueDatePicker/style/main.scss'
</style>
