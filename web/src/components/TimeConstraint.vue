<template>
  <div>
    <div v-if="Object.keys(date).length === 0 && date.constructor === Object" class="d-flex align-items-center">
      <CBadge style="font-size: 0.875rem;" color="primary">Forever</CBadge>
    </div>
    <div v-else class="mt-1" />

    <template v-for="(constraint, ctype) in date">
      <div
        v-for="(date_obj, index) in constraint"
        :key="ctype+'_'+index"
        style="white-space:pre"
      >
        <!-- Datetime -->
        <div v-if="ctype == 'datetime'" class="d-flex align-items-center mb-1">
          <CBadge style="font-size: 0.875rem;" color="info">
            {{ prettyDate(date_obj.from, false) }}
          </CBadge>
          <i class="la la-arrow-right la-lg"></i>
          <CBadge style="font-size: 0.875rem;" color="primary">
            {{ prettyDate(date_obj.until, false) }}
          </CBadge>
        </div>

        <!-- Time -->
        <div v-else-if="ctype == 'time'" class="d-flex align-items-center mb-1">
          <CBadge style="font-size: 0.875rem;" color="quaternary">
            {{ prettyDate(date_obj.from, false) }}
          </CBadge>
          <i class="la la-arrow-right la-lg"></i>
          <CBadge style="font-size: 0.875rem;" color="danger">
            {{ prettyDate(date_obj.until, false) }}
          </CBadge>
        </div>

        <!-- Weekdays -->
        <div v-else-if="ctype == 'weekdays'" class="d-flex align-items-center flex-wrap">
          <CBadge
            v-for="(weekday, ind) in date_obj.weekdays.sort()"
            :key="ind"
            :class="ind != date_obj.weekdays.length - 1 ? 'me-1 mb-1' : 'mb-1'"
            style="font-size: 0.875rem;"
            color="warning"
          >
            {{ getWeekday(weekday) }}
          </CBadge>
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { prettyDate, getWeekday } from '@/utils/functions'

export default defineComponent({
  props: {
    date: {type: Object, default: () => new Object()},
  },
  methods: {
    prettyDate,
    getWeekday,
  },
})
</script>
