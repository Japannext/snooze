import { reactive } from 'vue'

// Return a pagination reactive object compatible with
// naive-ui `n-data-table`
export function usePagination() {
  const pagination = reactive({
    page: 1,
    pageSize: 20,
    showSizePicker: true,
    pageSizes: [5, 10, 20, 30, 50, 100],
    onChange: (page: number) => {
      pagination.page = page
    },
    onUpdatePageSize: (pageSize: number) => {
      pagination.pageSize = pageSize
      pagination.page = 1
    }
  })
  return pagination
}
