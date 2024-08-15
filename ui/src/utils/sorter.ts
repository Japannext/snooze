import { reactive } from 'vue'

// Return a pagination reactive object compatible with
// naive-ui `n-data-table`
export function useSorter() {
  const sorter = reactive({
    sortName(order) {
      sorter.sortOrder = order
    },
    clearSorter() {
      sorter.sortOrder = false
    },
    handlerSorterChange(sorter) {
    }

  })
  return sorter
}
