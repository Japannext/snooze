<script setup lang="ts">
import { defineProps, computed } from 'vue'
import { DateTime } from 'luxon'
import { useTimeRelative } from '@/stores'

const { timeRelative } = useTimeRelative()

const props = defineProps<{
  value: number,
}>()

const date = computed(() => {
  return DateTime.fromMillis(props.value, {zone: "system"})
})

</script>

<template>
  <div>
    <span v-if="timeRelative">
      {{ date.toRelative() }}
    </span>
    <span v-else>
      {{ date.toFormat("yyyy-MM-dd HH:mm:ss") }} ({{ date.zoneName }})
    </span>
  </div>
</template>
