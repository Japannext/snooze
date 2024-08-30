<script setup lang="ts">
import { defineProps, computed } from 'vue'
import { DateTime } from 'luxon'
import { useTimeRelative } from '@/stores'

const { timeRelative } = useTimeRelative()

const props = defineProps<{
  timestampMillis: number,
}>()

const date = computed(() => {
  return DateTime.fromMillis(props.timestampMillis)
})

</script>

<template>
  <div>
    <span v-if="timeRelative">
      {{ date.toRelative() }}
    </span>
    <span v-else>
      {{ date.toFormat("yyyy-MM-dd HH:mm:ss") }}
    </span>
  </div>
</template>
