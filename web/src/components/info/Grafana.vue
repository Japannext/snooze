<template>
  <div>
    <CCard header="Grafana" header-class="text-center font-weight-bold">
      <CCardHeader class="text-center" style="font-weight:bold">
        Grafana
      </CCardHeader>
      <CCardBody class="p-2 m-auto">
        <a :href="data.rule_url">
          <img
            :src="data.image_url"
            class="rounded border border-danger"
            style="max-width:100%; max-height: 100%;"
          />
        </a>
      </CCardBody>
    </CCard>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'

interface GrafanaData {
  rule_url: string
  image_url: string
}

export default defineComponent({
  name: 'Grafana',
  props: {
    data: {type: Object as PropType<GrafanaData>, required: true},
  },
})
</script>
