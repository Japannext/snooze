<template>
  <div>
    <div v-for="(v, k) in metadata" :key="k">
      <CDropdownDivider v-if="k.includes('sweb_separator')" />
      <h4 v-else-if="k.includes('sweb_title')">
        <CDropdownDivider />{{ v }}<CDropdownDivider />
      </h4>
      <Base
        v-else
        :ref="`form.${k}`"
        v-model="dataValue[k]"
        :metadata="v"
        :data="modelValue"
        class="pb-2"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'

import Base from '@/components/form/Base.vue'

export default defineComponent({
  name: 'Form',
  components: { Base },
  props: {
    modelValue: {type: Object as PropType<object>, required: true},
    metadata: {type: Object as PropType<object>, default: () => new Object()},
  },
  emits: ['update:modelValue'],
  data() {
    return {
      dataValue: this.modelValue,
    }
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
})
</script>
