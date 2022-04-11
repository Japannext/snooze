<template>
  <div class="animated fadeIn">
    <List
      endpoint_prop="kv"
      edit_mode
      delete_mode
      add_mode
    >
    </List>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import List from '@/components/List.vue'

export default defineComponent({
  components: {
    List,
  },
})
</script>
