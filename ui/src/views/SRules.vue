<script setup lang="ts">
import { h } from 'vue'

import { RuleService } from '@/api'
import { STable, SRuleForm, SCondition } from '@/components'

const columns = [
  {title: "Name", key: "name"},
  {title: "Description", key: "description"},
  {title: "Condition", key: "condition", render: (row) => h(SCondition, {value: row.condition})},
]

</script>

<template>
  <div>
    <s-table
      name="Rule"
      :columns="columns"
      :service="RuleService"
      :form="SRuleForm"
      :crud="true"
    />
  </div>
</template>
