<script setup lang="ts">
import { ref, useSlots, withDefaults } from 'vue'
import { NButton } from 'naive-ui'

// Vue "Fallthrough attributes" will do the rest, since
// there is a single root element in the template
interface Props {
  onClick: () => Promise<void>,
  onFinish: () => void,
}

const props = withDefaults(defineProps<Props>(), {
  onFinish: () => {},
})

const loading = ref<boolean>(false)

function handleClick() {
  loading.value = true
  props.onClick()
  .then(() => {
    loading.value = false
    if (props.onFinish !== undefined) {
      props.onFinish()
    }
  })
}

const slots = useSlots()

</script>

<template>
  <n-button
    :loading="loading"
    @click="handleClick"
  >
    <!-- Passing parent slots -->
    <template v-for="(index, name) in slots" #[name]>
      <slot :name="name" />
    </template>
  </n-button>
</template>
