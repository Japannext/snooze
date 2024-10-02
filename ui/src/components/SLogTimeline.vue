<script setup lang="ts">
import {NSpace, NTimeline, NTimelineItem, NPagination} from 'naive-ui'

const items = ref<Array<TimelineItem>>([])

</script>

<template>
  <n-space vertical>
  <n-timeline
    
  >
    <n-timeline-item v-for="item in items">
      
    </n-timeline-item>
  </n-timeline>
  <n-pagination
    
  />
  </n-space>
</template>
