export type ListOf<T> = {
  items: Array<T>;
  total: number;
  more: boolean;
}

export type KeyValue = {
  [key: string]: string
}

export type Destination = {
  queue: string
  profile: string
}

export type NaiveColor = "default" | "error" | "primary" | "info" | "success" | "warning"
