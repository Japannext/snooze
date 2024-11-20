<script setup lang="ts">
import { defineProps, computed } from 'vue'
import { NPopover } from 'naive-ui'
import { DateTime } from 'luxon'
import { useTimeRelative } from '@/stores'

const { timeRelative } = useTimeRelative()

const props = defineProps<{
  ts: number,
  format?: string,
}>()

const date = computed(() => {
    return DateTime.fromMillis(props.ts, {zone: "system"})
})

const isRelative = computed<boolean>(() => {
  switch (props.format) {
    case "relative":
      return true
    case "absolute":
      return false
  }
  return timeRelative.value
})

</script>

<template>
  <div>
    <template v-if="ts">
      <span v-if="isRelative">
        {{ date.toRelative() }}
      </span>
      <n-popover v-else>
        <template #trigger>
          <span>
            {{ date.toFormat("yyyy-MM-dd HH:mm:ss") }}
          </span>
        </template>
        <span>Timezone: {{ date.zoneName }}</span>
      </n-popover>
    </template>
    <template v-else>
      Undefined
    </template>
  </div>
</template>
