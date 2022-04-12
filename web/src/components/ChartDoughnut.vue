<template>
  <SChart
    type="doughnut"
    :data="data"
    :options="default_options"
  />
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { gen_palette, hexToRgba } from '@/utils/colors'

import SChart from '@/components/SChart.vue'

export default defineComponent({
  name: 'ChartDoughnut',
  components: { SChart },
  props: {
    datasets: {type: Object, default: () => new Object()},
  },
  computed: {
    default_options () {
      return {
        maintainAspectRatio: false,
      }
    },
    data() {
      return {
        doughnut: {
          labels: [],
          datasets: this.datasets,
        },
      }
    },
    /**
      if (this.datasets) {
        var series_array = Object.keys(this.datasets).map(key => [key, this.datasets[key]])
        series_array = series_array.sort((a, b) => (a[1] < b[1]) ? 1: -1)
        var palette = gen_palette(series_array.length)
        return {
          datasets: [{
            borderColor: palette,
            backgroundColor: palette.map(x => hexToRgba(x, 75)),
            data: series_array.map(x => x[1]),
          }],
          labels: series_array.map(x => x[0]),
        }
      } else {
        return {datasets:[], labels: []}
      }
    },
  **/
  },
})
</script>
