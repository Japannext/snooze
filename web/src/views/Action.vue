<template>
  <div class="animated fadeIn">
    <List
      endpoint_prop="action"
      edit_mode
      delete_mode
      add_mode
    >
      <template #details_side="row">
        <AuditLogs collection="action" :object="row.item" />
      </template>
    </List>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import AuditLogs from '@/components/AuditLogs.vue'
import List from '@/components/List.vue'

export default defineComponent({
  components: {
    AuditLogs,
    List,
  },
})
</script>
