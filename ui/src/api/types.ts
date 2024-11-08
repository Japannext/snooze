import { DateTime } from 'luxon'

export type Log = {
  id?: string;

  source: Source;
  identity?: object
  labels?: object;

  actualTime: number;
  observedTime: number;
  displayTime: number;

  groupHash?: string;
  groupLabels?: any;

  severityText: string;
  severityNumber: number;

  message: string;

  mute?: Mute;
}

export type Snooze = {
  id?: string;
  groupName: string;
  hash: string;
  reason: string;
  startAt: DateTime;
  expireAt: DateTime;
}

export type ListOf<Type> = {
  items: Array<Type>;
  total: number;
}

export type LogResults = {
  items: Array<Log>;
  total: number;
}

export type Alert = {
}

export type AlertResults = {
  items: Array<Alert>;
  total: number;
  more: boolean;
}

export type Notification = {
  id?: string;

  notificationTime: number;

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
  queue: string
  profile: string
}

export type Mute = {
  skipNotification: boolean;
  skipStorage: boolean;
  reason: string;
}
