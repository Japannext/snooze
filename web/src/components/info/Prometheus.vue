<template>
  <div>
    <CCard header="Prometheus" header-class="text-center font-weight-bold">
      <CCardHeader class="text-center" style="font-weight:bold">
        Prometheus
      </CCardHeader>
      <CCardBody class="p-2">
        <div v-if="!!data.generatorURL">
          <a :href="data.generatorURL" target="_blank">Prometheus Graph</a>
        </div>
        <div v-if="!!data.externalURL" class="pt-1">
          External URL: <a :href="data.externalURL" target="_blank">{{ data.externalURL }}</a>
        </div>
      </CCardBody>
    </CCard>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

interface PrometheusData {
  generatorURL: string
  externalURL: string
}

export default defineComponent({
  name: 'Prometheus',
  props: {
    data: {type: Object},
  },
  data () {
    return {
    }
  },
})
</script>
