<template>
  <div>
    <SFormTags
      v-model="dataValue"
      :tags-options="items"
      :primary="primary"
      size="lg"
      :colorize="colorize"
      :required="required"
      trim
    >
    </SFormTags>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'

import { api2 } from '@/api2'
import { get_color } from '@/utils/colors'

import Base from './Base.vue'
import SFormTags from '@/components/SFormTags.vue'

export default defineComponent({
  components: {
    SFormTags,
  },
  extends: Base,
  props: {
    modelValue: {type: Array as PropType<string[]>, default: () => []},
    metadata: {type: Object, default: () => new Object()}
  },
  emits: ['update:modelValue'],
  data() {
    return {
      apiEndpoint: api2.endpoint(this.metadata.endpoint),
      primary: this.metadata.primary,
      colorize: this.metadata.colorize,
      required: this.metadata.required,
      dataValue: this.modelValue,
      importKeysData: this.importKeys,
      items: [] as object[],
    }
  },
  watch: {
    dataValue: {
      handler() {
        this.$emit('update:modelValue', this.dataValue)
      },
      immediate: true
    },
  },
  mounted () {
    this.reload()
  },
  methods: {
    get_color,
    reload() {
      this.apiEndpoint.find()
        .then((results: object[]) => {
          this.items = results
        })
    },
  },
})
</script>
