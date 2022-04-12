<template>
  <div>
    <CFormSelect
      v-model="dataValue.selected"
      :value="dataValue.selected"
      :invalid="isInvalid"
      @change="onChange"
    >
      <option disabled value="" :selected="selected == ''">Please select an option</option>
      <option v-for="item in items" :key="item[primary]" :value="item[primary]">
        {{ item['name'] }}
      </option>
    </CFormSelect>
    <CFormFeedback invalid>
      Field is required
    </CFormFeedback>
    <Form
      v-if="selection && selection[form]"
      v-model="dataValue.subcontent"
      :metadata="selection[form]"
      class="pt-2"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { api2 } from '@/api2'
import { DatabaseItem } from '@/utils/types'

import Base from '@/components/form/Base.vue'
import Form from '@/components/Form.vue'

interface ModelValue {
  selected: string
  subcontent: object
}

function defaultModelValue(): ModelValue {
  return {
    selected: '',
    subcontent: {},
  }
}

export default defineComponent({
  name: 'APIElement',
  components: {
    Form
  },
  extends: Base,

  props: {
    modelValue: {type: Object as PropType<ModelValue>, default: defaultModelValue},
    metadata: {type: Object, default: () => new Object()},
  },

  emits: ['update:modelValue'],

  data() {
    return {
      apiEndpoint: api2.endpoint(this.metadata.endpoint),
      primary: this.metadata.primary as string,
      form: this.metadata.form,
      required: this.metadata.required,
      dataValue: this.modelValue,
      items: [] as DatabaseItem[],
    }
  },

  computed: {
    isInvalid() {
      return (this.required && this.dataValue.selected == '')
    },
    selection() {
      return this.items.find((item: DatabaseItem) => {
        // We're making the type assumption that `item` has a key `this.primary`
        item[this.primary as keyof DatabaseItem] == this.dataValue.selected
      })
    },
  },

  watch: {
    dataValue: {
      handler() {
        this.$emit('update:modelValue', this.dataValue)
      },
      deep: true,
      immediate: true,
    },
  },
  mounted() {
    this.reload()
  },

  methods: {
    reload() {
      this.apiEndpoint.find()
        .then(results => {
          this.items = results
        })
    },
    onChange() {
      this.dataValue.subcontent = {}
    }
  },

})
</script>
