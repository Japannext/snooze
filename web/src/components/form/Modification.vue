<template>
  <div>
    <CForm @submit.prevent>
      <CRow v-for="(val, index) in dataValue" :key="index" class="m-0 pb-2">
        <CCard no-body class="p-0 col">
          <CCardBody class="p-0">
            <CInputGroup>
              <CFormSelect v-model="val[0]" :value="val[0]">
                <option v-for="opts in OPERATIONS" :key="opts.value" :value="opts.value">{{ opts.text }}</option>
              </CFormSelect>
              <template v-if="val[0] == 'DELETE'">
                <CFormInput id="field" v-model="val[1]" placeholder="Field to delete" />
              </template>
              <template v-else-if="val[0] == 'REGEX_PARSE'">
                <CFormInput id="field" v-model="val[1]" placeholder="Field to parse" />
              </template>
              <template v-else-if="val[0] == 'REGEX_SUB'">
                <CFormInput id="field" v-model="val[1]" placeholder="Field to parse" />
                <CFormInput id="out_field" v-model="val[2]" placeholder="Output field" />
              </template>
              <template v-else-if="val[0] == 'KV_SET'">
                <CFormInput id="dict" v-model="val[1]" placeholder="Dictionary" />
                <CFormInput id="field" v-model="val[2]" placeholder="Field" />
                <CFormInput id="out_field" v-model="val[3]" placeholder="Output field" />
              </template>
              <template v-else>
                <CFormInput id="field" v-model="val[1]" placeholder="Field" />
                <SFormInput id="value" v-model="val[2]" placeholder="Value" />
              </template>
              <CButton
                v-c-tooltip="{content: DOCUMENTATION[val[0]], trigger: 'click', placement: 'bottom'}"
                class="ms-auto"
                size="sm"
                color="secondary"
              >
                <i class="la la-info la-lg"></i>
              </CButton>
              <CButton
                size="sm"
                color="danger"
                @click="remove(index)"
                @click.stop.prevent
              >
                <i class="la la-trash la-lg"></i>
              </CButton>
            </CInputGroup>
            <CInputGroup v-if="val[0] == 'REGEX_PARSE'">
              <CFormTextarea id="regex" v-model="val[2]" placeholder="Regex with capture groups (?P<field_name>.*?)" />
            </CInputGroup>
            <CInputGroup v-else-if="val[0] == 'REGEX_SUB'">
              <CFormTextarea id="regex" v-model="val[3]" placeholder="Regex pattern to search for replacement" />
              <CFormTextarea id="sub" v-model="val[4]" placeholder="Substitute" />
            </CInputGroup>
          </CCardBody>
        </CCard>
      </CRow>
    </CForm>
    <CCol xs="auto">
      <CButton
        v-c-tooltip="'Add'"
        color="secondary"
        @click="append"
        @click.stop.prevent
      >
        <i class="la la-plus la-lg"></i>
      </CButton>
    </CCol>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Base from './Base.vue'
import SFormInput from '@/components/SFormInput.vue'

const DOCUMENTATION = {
  'SET': "Set a field to a given value (string)",
  'DELETE': "Delete a field value",
  'ARRAY_APPEND': "Append a string to an array",
  'ARRAY_DELETE': "Delete an element from an array by value",
  'REGEX_PARSE': "Given a regex with named capture groups, the value of the capture groups will be merged to the record by name",
  'REGEX_SUB': "Search the elements matching a regex, and replace them with a substitute",
  'KV_SET': "Map a field to a value in a key-value dictionary",
}
const OPERATIONS = [
  {value: 'SET', text: 'Set'},
  {value: 'DELETE', text: 'Delete'},
  {value: 'ARRAY_APPEND', text: 'Append (to array)'},
  {value: 'ARRAY_DELETE', text: 'Delete (from array)'},
  {value: 'REGEX_PARSE', text: 'Regex parse (capture)'},
  {value: 'REGEX_SUB', text: 'Regex sub'},
  {value: 'KV_SET', text: 'Key-value mapping'},
]

export default defineComponent({
  name: 'Modification',
  components: {
    SFormInput,
  },
  extends: Base,
  props: {
    modelValue: {type: Array, default: () => []},
  },
  emits: ['update:modelValue'],
  data () {
    return {
      OPERATIONS: OPERATIONS,
      DOCUMENTATION: DOCUMENTATION,
      dataValue: this.modelValue,
    }
  },
  watch: {
    dataValue: {
      handler: function () {
        this.$emit('update:modelValue', this.dataValue)
      },
      immediate: true
    },
  },
  methods: {
    append() {
      this.dataValue.push(['SET', '', ''])
    },
    remove(index: number) {
      this.dataValue.splice(index, 1)
    },
  },
})
</script>

<style scoped lang="scss">
.input-group {
  .form-select, .form-control, btn {
    margin: -1px;
  }
}
</style>
