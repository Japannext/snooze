import { reactive, ref, type VNodeChild } from 'vue'
import { useRouteQuery } from '@vueuse/router'
import type { PaginationProps } from 'naive-ui'

/*
export interface Pagination {
  page: number
  pageSize: number
  showSizePicker: boolean
  pageSizes: number[]
  itemCount: number
  prefix(p: Pagination): void
  onUpdatePage(p: number): void
  onUpdatePageSize(p: number): void
  setMore(m: boolean): void
}
*/

// Pagination set of parameters for the API calls
export interface Pagination {
  page?: number
  pageSize?: number
}

// A naive-ui compatible pagination that can be used with n-data-table
export interface NaivePagination extends PaginationProps {
  setMore(m: boolean): void
}

type PrefixOption = {
  startIndex: number,
  endIndex: number,
  pageSize: number,
  pageCount: number,
  itemCount: number|undefined
}

// Return a pagination reactive object compatible with
// naive-ui `n-data-table`
export function usePagination(refresh: Function): NaivePagination {
  const page = useRouteQuery<number>('page', 1, {transform: Number})
  const size = useRouteQuery<number>('size', 20, {transform: Number})
  const more = ref<boolean>(false)
  const pagination = reactive({
    page: page,
    pageSize: size,
    showSizePicker: true,
    pageSizes: [5, 10, 20, 30, 50, 100],
    itemCount: 0,
    prefix(p: PrefixOption): VNodeChild {
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
