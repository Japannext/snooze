<template>
  <div>
    <VueDatePicker
      v-model="dataValue"
      time-picker
      placeholder="Select time"
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
import { TimeConstraint } from '@/utils/types'

import VueDatePicker from '@vuepic/vue-datepicker';
import Base from './Base.vue'

const TIME_FORMAT = "HH:mm:ssZ"

function defaultTimeConstraint(): TimeConstraint {
  const now = moment()
  const oneHourLater = now.clone().add(1, 'hours')
  return {
    from: now.format(TIME_FORMAT),
    until: oneHourLater.format(TIME_FORMAT),
  }
}

interface PickerTime {
  hours: number
  minutes: number
}

function pickerToConstraint(data: [PickerTime|null, PickerTime|null]): TimeConstraint {
  const [from, until] = data
  let constraint: TimeConstraint = {}
  if (from !== null) {
    const fromText = `${from.hours}:${from.minutes}`
    constraint.from = moment(fromText, 'HH:mm').format(TIME_FORMAT)
  }
  if (until !== null) {
    const untilText = `${until.hours}:${until.minutes}`
    constraint.until = moment(untilText, 'HH:mm').format(TIME_FORMAT)
  }
  return constraint
}

function constraintToPicker(data: TimeConstraint): [PickerTime|null, PickerTime|null] {
  const from = data.from
  const until = data.until
  let fromPicker = null
  let untilPicker = null
  if (from !== undefined) {
    const fromDate = moment(from)
    fromPicker = {hours: fromDate.hours(), minutes: fromDate.minutes()}
  }
  if (until !== undefined) {
    const untilDate = moment(until)
    untilPicker = {hours: untilDate.hours(), minutes: untilDate.minutes()}
  }
  return [fromPicker, untilPicker]
}

export default defineComponent({
  name: 'TimeConstraint',
  components: { VueDatePicker },
  extends: Base,
  props: {
    modelValue: {type: Object as PropType<TimeConstraint>, default: defaultTimeConstraint},
  },
  emits: ['update:modelValue'],
  data() {
    return {
      dataValue: constraintToPicker(this.modelValue),
    }
  },
  watch: {
    dataValue: {
      handler() {
        this.$emit('update:modelValue', pickerToConstraint(this.dataValue))
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
