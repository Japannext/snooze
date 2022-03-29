<template>
  <CCard>
    <CCardHeader class='text-center' style='font-weight:bold'>
      Error in plugin {{ corePlugin }}
    </CCardHeader>
    <CCardBody class="p-0">
      <h5 class="ps-2 pt-1">{{ errorClass }}: {{ message }}</h5>
      <p>
        {{ tracebackArray[0] }}
        <span v-if="more">{{ tracebackArray[1] }}</span>
      </p>
    </CCardBody>
  </CCard>
</template>

<script>
import { more } from '@/utils/api'

export default {
  name: 'Error',
  components: {
  },
  props: {
    data: {type: Object},
    errorClass: {type: String},
    message: {type: String},
    traceback: {type: String},
  },
  data () {
    return {
      more: false,
      corePlugin: this.data.core_plugin,
      errorClass: this.data.error_class,
      message: this.data.message,
      traceback: this.data.traceback,
    }
  },
  computed: {
    tracebackArray () {
      if (this.traceback) {
        return more(this.trackback)
      } else {
        return false
      }
    },
  },
}
</script>
