<script setup lang="ts">

import { ref, h, reactive, onMounted, VNodeChild } from 'vue'

import {
  NTag, NIcon,
  NTime,
  NSpace,
  NTabs, NTab,
  NButton, NInputGroup
} from 'naive-ui'
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import { Server, Boat } from '@vicons/ionicons5'
import { Bolt, BellSlash } from '@vicons/fa'

import { AlertService } from '@/api'
import type { AlertV2, AlertV2Collection } from '@/api'

import STable from '@/components/STable.vue'
import SResource from '@/components/SResource.vue'
import SSearch from '@/components/SSearch.vue'
import SSeverity from '@/components/SSeverity.vue'
import SAttributes from '@/components/SAttributes.vue'

function renderTimestamp(item: object): VNodeChild {
  const date = new Date(item.timestamp)
  return h(NTime, {time: date})
}

const columns = [
  {
    title: "Timestamp",
    key: "timestamp",
    render: renderTimestamp,
    width: 160,
  },
  {
    title: "Resource",
    key: "resource",
    resizable: true,
    width: 300,
    render: (item) => h(SResource, {resource: item.resource, size: "small"}),
  },
  {
    title: "Body",
    key: "body",
  },
  {
    title: 'Hits',
    key: 'hits',
    width: 60,
  },
  {
    title: "Attributes",
    key: "attributes",
    resizable: true,
    render: (item) => h(SAttributes, {attributes: item.attributes, size: "small"}),
  },
  {
    title: "Severity",
    key: "severity.text",
    width: 100,
    render: (item) => h(SSeverity, {severity: item.severity, size: "small"}),
  },
]

</script>

<template>
  <div>
    <s-table
      :columns="columns"
      :service="AlertService"
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
