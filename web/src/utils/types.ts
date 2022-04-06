
export type QueryElement = string | Query[]
export type Query = Array<QueryElement>

export interface PaginationOptions {
  perpage?: number
  pagenb?: number
  asc?: boolean
  orderby?: string
}

export interface DatabaseItem {
  uid: string
}

export interface RejectedItem {
  error: string
}

export interface ChangeResult {
  added?: DatabaseItem[]
  updated?: DatabaseItem[]
  replaced?: DatabaseItem[]
  rejected?: RejectedItem[]
}

export interface Result<U> {
  data: U
  count: number
}

enum KeyDiffType {
  Added = 'added',
  Removed = 'removed',
  Updated = 'updated',
}

export type AuditSummary = {
  [key: string]: KeyDiffType,
}

export type SummaryEntry = {
  symbol: string,
  color: string,
  name: string,
}

export interface AuditMetadata {
  name: string
  color: string
  icon: string
  methodColor: string
  quickSummary: SummaryEntry,
  summaryCount: number
}

export interface AuditItem extends DatabaseItem {
  collection: string
  object_id: string
  timestamp: string
  action: 'added' | 'updated' | 'replaced' | 'deleted' | 'rejected'
  username: string
  method: string
  snapshot: object
  source_ip?: string
  user_agent?: string
  summary?: AuditSummary,
}

