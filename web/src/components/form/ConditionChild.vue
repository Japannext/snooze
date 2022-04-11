<template>
  <CForm @submit.prevent>
    <!-- No condition case (always true) -->
    <template v-if="dataValue.type == 'alwaysTrue'">
      <CBadge style="font-size: 0.875rem;" color="danger" class="col me-2 align-middle">
        Always true
      </CBadge>
      <CButton
        class="col"
        color="secondary"
        @click="dataValue = defaultCondition()"
      >
        <i class="la la-plus la-lg"></i>
      </CButton>
    </template>

    <!-- AND/OR/NOT case -->
    <template v-if="['logic', 'not'].includes(dataValue.type)">
      <CForm inline class="row g-0">
        <CCol xs="auto">
          <CInputGroup>
            <CFormSelect v-model="operation" :value="dataValue.operation" class="w-auto">
              <option v-for="option in LOGIC_OPTIONS" :key="option.value" :value="option.value">{{ option.text }}</option>
            </CFormSelect>
            <CButton v-if="dataValue.operation != 'NOT'" color="secondary" @click="logicAdd" @click.stop.prevent>
              <i class="la la-plus la-lg"></i>
            </CButton>
            <CButton color="danger" @click="escalateDelete" @click.stop.prevent>
              <i class="la la-trash la-lg"></i>
            </CButton>
          </CInputGroup>
        </CCol>
      </CForm>
      <ul class="ps-3">
        <CForm v-for="(arg, i) in dataValue.args" :key="arg.id" class="pt-1" @submit.prevent>
          <ConditionChild v-model="dataValue.args[i]" :index="i" @delete-event="deleteCondition" />
        </CForm>
      </ul>
    </template>

    <!-- `field=value` case (and other operations) -->
    <template v-else-if="['binary', 'unary'].includes(dataValue.type)">
      <CInputGroup>
        <CFormInput v-model="dataValue.args[0]" placeholder="Field" style="flex: 0 0 auto; width: 25%" />
        <CFormSelect v-model="operation" :value="dataValue.operation" style="flex: 0 0 auto; width: 15%">
          <option v-for="option in OPERATION_OPTIONS" :key="option.value" :value="option.value">{{ option.text }}</option>
        </CFormSelect>
        <SFormInput v-if="dataValue.type == 'binary'" v-model="dataValue.args[1]" placeholder="Value" />
        <CButton color="secondary" @click="fork" @click.stop.prevent>
          <i class="la la-plus la-lg"></i>
        </CButton>
        <CButton color="info" @click="dataValue = defaultCondition()" @click.stop.prevent>
          <i class="la la-redo-alt la-lg"></i>
        </CButton>
        <CButton color="danger" @click="escalateDelete" @click.stop.prevent>
          <i class="la la-trash la-lg"></i>
        </CButton>
      </CInputGroup>
    </template>
  </CForm>
</template>

<script lang="ts">

// Options for the operation selector
const OPERATION_OPTIONS = [
  {value: '=', text: '='},
  {value: '!=', text: '!='},
  {value: '>', text: '>'},
  {value: '>=', text: '>='},
  {value: '<', text: '<'},
  {value: '<=', text: '<='},
  {value: 'MATCHES', text: 'matches'},
  {value: 'EXISTS', text: 'exists?'},
  {value: 'CONTAINS', text: 'contains'},
  {value: 'SEARCH', text: 'search'},
]

// Options for the logic selector
const LOGIC_OPTIONS = [
  {value: 'OR', text: 'OR'},
  {value: 'AND', text: 'AND'},
  {value: 'NOT', text: 'NOT'},
]

import { defineComponent } from 'vue'

import { ConditionObject, OPERATION_TYPE } from '@/utils/condition2'
import SFormInput from '@/components/SFormInput.vue'

// Return a new condition object
function defaultCondition(): ConditionObject {
  return new ConditionObject('=', ['', ''])
}

export default defineComponent({
  name: 'ConditionChild',
  components: {
    SFormInput,
  },
  props: {
    modelValue: {type: ConditionObject, required: true},
    index: {type: Number, default: 0},
    root: {type: Boolean, default: false},
  },
  emits: ['update:modelValue', 'delete-event'],
  data () {
    return {
      OPERATION_TYPE: OPERATION_TYPE,
      OPERATION_OPTIONS: OPERATION_OPTIONS,
      LOGIC_OPTIONS: LOGIC_OPTIONS,
      dataValue: this.modelValue as ConditionObject,
    }
  },
  computed: {
    // We compute the operation for logic selectors in order to correct the number
    // of arguments that is different for OR/AND and NOT.
    operation: {
      get() { return this.dataValue.operation },
      set(op: string) {
        if (op == 'NOT') {
          const cond = this.dataValue.args[0]
          if (cond instanceof ConditionObject) {
            this.dataValue.args = [cond]
          } else {
            throw `Unexpected argument type for NOT: ${cond} (${typeof cond})`
          }
        } else if ((op == 'AND' || op == 'OR') && this.dataValue.args.length < 2) {
          this.logicAdd()
        }
        this.dataValue.operation = op
      },
    },
  },
  watch: {
    dataValue: {
      handler() {
        this.$emit('update:modelValue', this.dataValue)
      },
    }
  },
  methods: {
    defaultCondition,

    // Triggered when we push the `+` button for a logic operator (AND/OR)
    logicAdd() {
      this.dataValue.args.push(this.defaultCondition())
    },
    // Delete a condition at the given index in a AND/OR condition
    deleteCondition(index: number) {
      this.dataValue.args.splice(index, 1)
      // Correct the condition if it's invalid (AND/OR less than 2 argument,
      // NOT with less than 1 argument)
      if (this.dataValue.type == 'logic' && this.dataValue.args.length < 2) {
        this.dataValue = this.dataValue.args[0] as ConditionObject
      }
      if (this.dataValue.type == 'not' && this.dataValue.args.length < 1) {
        this.dataValue = this.defaultCondition()
      }
    },
    // Triggered when the delete button is pressed for any condition. This will
    // escalate the delete operation to the parent condition, or reset the condition
    // if it's the root condition.
    escalateDelete() {
      if (this.root) {
        this.dataValue = this.defaultCondition()
      } else {
        this.$emit('delete-event', this.index)
      }
    },
    // Trigerred when pushing the `+` button for a normal condition (a=x)
    // This will create a logic operator at the place of the condition, resulting
    // in "a=1 AND defaultCondition()"
    fork() {
      this.dataValue = this.dataValue.combine('AND', this.defaultCondition())
    },
  },
})
</script>
