import { reactive, ref } from 'vue'

import { useRouteQuery } from '@vueuse/router'

// Return a pagination reactive object compatible with
// naive-ui `n-data-table`
export function usePagination(refresh: Function) {
  const page = useRouteQuery('page', 1, {transform: Number})
  const size = useRouteQuery('size', 20, {transform: Number})
  const more = ref<boolean>(false)
  const pagination = reactive({
    page: page,
    pageSize: size,
    showSizePicker: true,
    pageSizes: [5, 10, 20, 30, 50, 100],
    itemCount: 0,
    prefix(p) {
      console.log(p)
      if (more.value) {
        return `${p.itemCount}+ objects`
      } else {
        return `${p.itemCount} objects`
      }
    },
    onUpdatePage(currentPage: number) {
      page.value = currentPage
      refresh()
    },
    onUpdatePageSize(pageSize: number) {
      size.value = pageSize
      page.value = 1
      refresh()
    },
    setMore(m: boolean) {
      more.value = m
    },
  })
  return pagination
}
