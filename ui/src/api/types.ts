export type Log = {
  source: Source;

  timestamp?: number;
  observedTimestamp?: number;

  groupHash?: string;
  groupLabels?: any;

  severityText: string;
  severityNumber: number;

  labels?: any;
  attributes?: any;
  body?: any;

  mute?: Mute;
}

export type LogsResponse = {
  logs: Array<Log>;
}

export type Source = {
  kind: string
  name?: string
}

export type Mute = {
  enabled: boolean;
  component: string;
  rule: string;
  skipNotification: boolean;
  skipStorage: boolean;
  silentTest: boolean;
  humanTest: boolean;
}

