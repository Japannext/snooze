<script setup lang="ts">

import { ref, h, reactive, onMounted, VNodeChild } from 'vue'

import {
  NTag, NIcon, NSpan,
  NTime,
  NSpace,
  NTabs, NTab,
  NButton, NInputGroup
} from 'naive-ui'
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import { Bolt, BellSlash } from '@vicons/fa'

import { AlertgroupService } from '@/api'
import type { AlertV2, AlertV2Collection } from '@/api'

import STable from '@/components/STable.vue'
import SGroup from '@/components/SGroup.vue'
import SHash from '@/components/SHash.vue'
import SSearch from '@/components/SSearch.vue'
import SSeverity from '@/components/SSeverity.vue'
import SAttributes from '@/components/SAttributes.vue'

function renderTimestamp(item: object): VNodeChild {
  const timestamp = item.last_hit
  try {
    const date = new Date(timestamp)
    return h(NTime, {time: date})
  } catch (error) {
    console.log(error)
    return h(NSpace, timestamp)
  }
}

const columns = [
  {
    title: "Last hit",
    key: "last hit",
    render: renderTimestamp,
    width: 160,
  },
  {
    title: 'Hits',
    key: 'hits',
    width: 55,
  },
  {
    title: "Group",
    key: "group",
    resizable: true,
    width: 300,
    render: (item) => h(SGroup, {group: item.group, size: "small"}),
  },
  {
    title: "Last message",
    key: "last_message",
  },
]

</script>

<template>
  <div>
    <s-table
      :columns="columns"
      :service="AlertgroupService"
      :crud="false"
    >
      <template #left>
        <n-tabs default-value="Alerts">
          <n-tab name="Alerts">
            <n-icon :component="Bolt" />
            Alerts
          </n-tab>
          <n-tab name="Snoozed" align="end">
            <n-icon :component="BellSlash" />
            <div>Snoozed</div>
          </n-tab>
        </n-tabs>
      </template>
    </s-table>
  </div>
</template>
