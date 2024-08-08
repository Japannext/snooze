import { ref } from 'vue'
import { useLocalStorage } from '@vueuse/core'

function useUi() {

  const devMode = ref<boolean>(false)
  const sideCollapsed = ref<boolean>(false)

}
