<template>
  <div v-for="widget in widgets" :key="widget.name" class="d-inline-flex m-auto pe-2">
    <component
      :is="widget.vue_component"
      v-if="widget.vue_component"
      :id="`component_${widget.vue_component}`"
      :options="widget"
      class="pb-1 m-auto"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { api2 } from '@/api2'
import { DatabaseItem } from '@/utils/types'

import PatliteWidget from '@/components/PatliteWidget.vue'

interface Widget extends DatabaseItem {
  enabled?: boolean
  vue_component?: string
}

// Create a card fed by an API endpoint.
export default defineComponent({
  name: 'WidgetList',
  components: {
    PatliteWidget,
  },
  props: {
  },
  data() {
    return {
      widgetEndpoint: api2.endpoint('widget'),
      widgets: [] as Widget[],
    }
  },
  mounted () {
    this.listWidgets()
  },
  methods: {
    /** Get the list of widgets from the API.
     * Update the `widgets` variable if the API return a result.
    **/
    listWidgets() {
      this.widgetEndpoint.find()
      .then(results => {
        const widgetResults: Widget[] = results
        this.widgets = widgetResults.filter((widget: Widget) => {
          widget.enabled === undefined || widget.enabled || widget.vue_component !== undefined
        })
      })
    },
  },
})
</script>
