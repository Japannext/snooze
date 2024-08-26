export type Log = {
  id?: string;

  source: Source;
  identity?: object
  labels?: object;

  timestampMillis: number;
  observedMillis: number;

  groupHash?: string;
  groupLabels?: any;

  severityText: string;
  severityNumber: number;

  message: string;

  mute?: Mute;
}

export type LogResults = {
  items: Array<Log>;
  total: number;
}

export type Notification = {
  id?: string;
  timestampMillis: number;
  destination: Destination;
  Acknowledged: boolean;
  alertUID?: string;
  logUID?: string;
  body: object;
}

export type NotificationResults = {
  items: Array<Notification>;
  total: number;
}

export type Source = {
  kind: string
  name?: string
}

export type Destination = {
  kind: string
  name: string
}

export type Mute = {
  skipNotification: boolean;
  skipStorage: boolean;
  reason: string;
}
