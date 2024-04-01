<script setup lang="ts">
import { h } from 'vue'
import { NTabs, NTab } from 'naive-ui'

import { RoleService } from '@/api'
import { STable, SRoleForm, SSet } from '@/components'

const columns = [
  {title: "Name", key: "name"},
  {title: "Description", key: "description"},
  {title: "Permissions", key: "scopes", render: (row) => h(SSet, {value: row.scopes})},
]

</script>

<template>
  <s-table
    name="Role"
    :columns="columns"
    :service="RoleService"
    :form="SRoleForm"
    :crud="true"
  >
    <template #left>
      <n-tabs default-value="Alerts">
        <n-tab name="Roles">
          <n-icon :component="Bolt" />
          Alerts
        </n-tab>
        <n-tab name="Role bindings" align="end">
          <n-icon :component="BellSlash" />
          <div>Role bindings</div>
        </n-tab>
      </n-tabs>
    </template>
  </s-table>
</template>
