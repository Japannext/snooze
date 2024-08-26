<script setup lang="ts">
import { defineProps } from 'vue'
import { DateTime } from 'luxon'
import { useTimeRelative } from '@/stores'

const { timeRelative } = useTimeRelative()

const props = defineProps<{
  timestampMillis: number,
}>()

</script>

<template>
  <div>
    <span v-if="timeRelative">
      {{ DateTime.fromMillis(props.timestampMillis).toRelative() }}
    </span>
    <span v-else>
      {{ DateTime.fromMillis(props.timestampMillis).toISO() }}
    </span>
  </div>
</template>
