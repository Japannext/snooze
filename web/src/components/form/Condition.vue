<template>
  <div>
    <!-- Condition dataValue: {{ dataValue }}, modelValue: {{ modelValue }} -->
    <ConditionChild
      v-model="dataValue"
      :root="true"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import { ConditionObject } from '@/utils/condition2'

import ConditionChild from '@/components/form/ConditionChild.vue'
import Base from './Base.vue'

// Convert the array input input a ConditionObject type to
// pass to the ConditionChild component
export default defineComponent({
  name: 'Condition',
  components: { ConditionChild },
  extends: Base,
  props: {
    modelValue: {type: Array, default: () => [""]},
  },
  emits: ['update:modelValue'],
  data () {
    return {
      dataValue: ConditionObject.fromArray(this.modelValue),
    }
  },
  watch: {
    dataValue: {
      handler() {
        this.$emit('update:modelValue', this.dataValue.toArray())
      },
      deep: true,
    }
  },
})
</script>
