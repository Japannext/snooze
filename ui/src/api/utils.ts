export type ListOf<T> = {
  items: Array<T>;
  total: number;
  more: boolean;
}

export type Destination = {
  queue: string
  profile: string
}
