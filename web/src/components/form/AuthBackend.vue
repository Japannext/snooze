<template>
  <div>
    <APIElement
      v-model="dataValue"
      endpoint="login"
      primary="endpoint"
      subkey="backends"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import Base from './Base.vue'
import APIElement from './APIElement.vue'

export default defineComponent({
  components: {
    APIElement,
  },
  extends: Base,
  props: {
    modelValue: {type: Object, default: () => new Object()},
  },
  data() {
    return {
      dataValue: this.modelValue,
    }
  },
})
</script>
