<template>
  <div class="h-100">
    <CForm
      class="d-flex align-items-center mb-0 h-100"
    >
      <CFormCheck
        v-for="weekday in WEEKDAYS"
        :key="weekday.value"
        v-model="dataValue[weekday.name]"
        inline
        :label="weekday.text"
      />
    </CForm>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { WeekdayConstraint, Week, Weekday, WeekdayNumber } from '@/utils/types'

import Base from './Base.vue'

function defaultWeekdayConstraint(): WeekdayConstraint {
  return {'weekdays': []}
}

const WEEKDAYS = [
  {name: 'sunday', text: 'Sunday', value: 0},
  {name: 'monday', text: 'Monday', value: 1},
  {name: 'tuesday', text: 'Tuesday', value: 2},
  {name: 'wednesday', text: 'Wednesday', value: 3},
  {name: 'thursday', text: 'Thursday', value: 4},
  {name: 'friday', text: 'Friday', value: 5},
  {name: 'saturday', text: 'Saturday', value: 6},
]

const WEEKDAYS_TO_NUMBERS = new Map<Weekday, WeekdayNumber>(
  WEEKDAYS.map(weekday => [weekday.name as Weekday, weekday.value as WeekdayNumber])
)
const NUMBERS_TO_WEEKDAYS = new Map<WeekdayNumber, Weekday>(
  WEEKDAYS.map(weekday => [weekday.value as WeekdayNumber, weekday.name as Weekday])
)

function weekToWeekdayConstraint(week: Week): WeekdayConstraint {
  let weekdays: Array<WeekdayNumber> = []
  for (let [key, value] of WEEKDAYS_TO_NUMBERS) {
    if (week.get(key)) {
      weekdays.push(value)
    }
  }
  return {weekdays: weekdays}
}

function weekdayConstraintToWeek(constraint: WeekdayConstraint): Week {
  let week: Week = new Map()
  constraint.weekdays.forEach((value: WeekdayNumber) => {
    const dayName = NUMBERS_TO_WEEKDAYS.get(value)
    if (dayName !== undefined) {
      week.set(dayName, true)
    }
  })
  return week
}

export default defineComponent({
  name: 'WeekdayConstraint',
  extends: Base,
  props: {
    modelValue: {type: Object as PropType<WeekdayConstraint>, default: defaultWeekdayConstraint},
  },
  emits: ['update:modelValue'],
  data() {
    return {
      WEEKDAYS: WEEKDAYS,
      dataValue: weekdayConstraintToWeek(this.modelValue),
    }
  },
  watch: {
    dataValue: {
      handler() { this.$emit('update:modelValue', weekToWeekdayConstraint(this.dataValue)) },
      immediate: true,
      deep: true,
    },
  },
})
</script>
