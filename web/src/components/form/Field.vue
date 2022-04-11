<template>
  <div>
    <SFormTags
      v-model="dataValue"
      placeholder="Enter fields separated by space"
      trim
    ></SFormTags>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import Base from '@/components/form/Base.vue'
import SFormTags from '@/components/SFormTags.vue'

export default defineComponent({
  name: 'Condition',
  components: {
    SFormTags,
  },
  extends: Base,
  props: {
    modelValue: {type: Array as PropType<string[]>, default: () => [] as string[]},
    metadata: {type: Object, default: new Object()},
  },
  emits: ['update:modelValue'],
  data () {
    return {
      dataValue: this.modelValue,
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
})
</script>
