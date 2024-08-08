export type Pagination = {
  index: number;
  perPage: number;
  orderBy?: string;
  direction?: number;
}

export function newPagination(): Pagination {
  return {index: 0, perPage: 20}
}
