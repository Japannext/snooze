<template>
  <div>
    <h5 v-if="Object.keys(dataValue).length === 0 && dataValue.constructor === Object">
      <CBadge color="primary">Forever</CBadge>
    </h5>
    <CForm class="m-0" @submit.prevent>
      <template v-for="(constraints, type) in dataValue">
        <CRow v-for="(value, index) in constraints" :key="`${type}-${index}`" class="mb-2 g-0">
          <CCol xs="1">
            <CButton color="danger" size="lg" !click="remove_component(constraint, index)" @click.stop.prevent>
              X
            </CButton>
          </CCol>
          <CCol xs="11" class="m-auto">
            <component
              :is="TYPE_TO_COMPONENT[type]"
              v-model="constraints[index]"
            />
          </CCol>
        </CRow>
      </template>
    </CForm>
    <CForm inline>
      <CRow class="g-0">
        <CCol xs="auto">
          <CInputGroup>
            <CFormSelect v-model="selectedComponent" class="col-form-label">
              <option v-for="opt in TIME_CONSTRAINTS" :key="opt.component" :value="opt.component">{{ opt.text }}</option>
            </CFormSelect>
            <CButton color="secondary" @click="add_component(selected)" @click.stop.prevent>
              <i class="la la-plus la-lg"></i>
            </CButton>
          </CInputGroup>
        </CCol>
      </CRow>
    </CForm>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import {
  TimeConstraints as TimeConstraintsType,
  TimeConstraint as TimeConstraintType,
  DatetimeConstraint as DatetimeConstraintType,
  WeekdayConstraint as WeekdayConstraintType
} from '@/utils/types'

import Base from './Base.vue'
import DatetimeConstraint from '@/components/form/DatetimeConstraint.vue'
import TimeConstraint from '@/components/form/TimeConstraint.vue'
import WeekdayConstraint from '@/components/form/WeekdayConstraint.vue'

const TIME_CONSTRAINTS = [
  {type: 'datetime', text: 'DateTime', component: 'DatetimeConstraint'},
  {type: 'time',     text: 'Time',     component: 'TimeConstraint'},
  {type: 'weekdays', text: 'Weekdays', component: 'WeekdayConstraint'},
]

const TYPE_TO_COMPONENT = new Map(TIME_CONSTRAINTS.map(el => [el.type, el.component]))

type ConstraintType = 'datetime'|'time'|'weekdays'

function addNew(constraints: any, key: ConstraintType) {
  switch(key) {
    case 'datetime': {
      (constraints as Array<DatetimeConstraintType>).push({})
      break
    }
    case 'time': {
      (constraints as Array<TimeConstraintType>).push({})
      break
    }
    case 'weekdays': {
      (constraints as Array<WeekdayConstraintType>).push({weekdays: []})
    }
  }
}


export default defineComponent({
  name: 'TimeConstraints',
  components: { DatetimeConstraint, TimeConstraint, WeekdayConstraint },
  extends: Base,
  props: {
    modelValue: {type: Object as PropType<TimeConstraintsType>, default: () => new Object()},
  },
  emits: ['update:modelValue'],
  data() {
    return {
      TYPE_TO_COMPONENT: TYPE_TO_COMPONENT,
      selectedComponent: 'datetime',
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
  methods: {
    addComponent(key: ConstraintType) {
      let constraints = this.dataValue[key]
      if (constraints === undefined) {
        constraints = []
      }
      addNew(constraints, key)
    },
    removeComponent(key: ConstraintType, index: number) {
      let constraints = this.dataValue[key]
      if (constraints !== undefined) {
        constraints.splice(index, 1)
        if (constraints.length == 0) {
          delete this.dataValue[key]
        }
      }
    },
  },
})
</script>
