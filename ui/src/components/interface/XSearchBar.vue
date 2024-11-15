<script setup lang="ts">
import { ref, defineModel, defineEmits } from 'vue'
import { NSpace, NInputGroup, NButton } from 'naive-ui'
import { XSearch, XTimeRange, XFilter } from '@/components/interface'

export type SearchParams = {
  filter: string;
  search: string;
  start?: number;
  end?: number;
}

const params = ref<SearchParams>({})

const emit = defineEmits<{
  (e: 'search', params: SearchParams): void
}>()
</script>

<template>
  <n-space style="padding: 5px; margin-bottom: 10px;">
    <x-filter v-model:value="params.filter" :filters="filters" @change="getItems" />
    <n-input-group>
      <x-time-range ref="stimerange" v-model:rangeMillis="rangeMillis" @update="onUpdateTimerange" />
      <x-search v-model:value="params.search" @search="onSearch" />
      <n-button @click="getItems()"><n-icon :component="Refresh" /></n-button>
    </n-input-group>
    <n-input-group>
      <template v-if="selectedItems.length > 0">
        <!-- Ack button -->
        <n-button type="primary" @click="showAckModal = true">Ack ({{ selectedItems.length }})</n-button>
        <x-ack-modal v-model:show="showAckModal" :ids="selectedItems" @success="unselect" />

        <n-button type="warning">Escalate ({{ selectedItems.length }})</n-button>
      </template>
    </n-input-group>
  </n-space>
</template>
