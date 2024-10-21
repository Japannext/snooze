<script setup lang="ts">
import { defineProps, computed } from 'vue'
import { DateTime } from 'luxon'
import { useTimeRelative } from '@/stores'
import type { Timestamp } from '@/api/types'

const { timeRelative } = useTimeRelative()

const props = defineProps<{
  ts: Timestamp,
}>()

const date = computed(() => {
    return DateTime.fromMillis(props.ts.display, {zone: "system"})
})

</script>

<template>
  <div>
    <template v-if="ts">
      <span v-if="timeRelative">
        {{ date.toRelative() }}
      </span>
      <span v-else>
        {{ date.toFormat("yyyy-MM-dd HH:mm:ss") }} ({{ date.zoneName }})
      </span>
    </template>
    <template v-else>
      Undefined
    </template>
  </div>
</template>
