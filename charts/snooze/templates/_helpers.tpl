{{- define "snooze.env" }}
- name: OPENSEARCH_USERNAME
  valueFrom:
    secretKeyRef:
      name: "{{ .Release.Name }}-opensearch-admin-password"
      key: username
- name: OPENSEARCH_PASSWORD
  valueFrom:
    secretKeyRef:
      name: "{{ .Release.Name }}-opensearch-admin-password"
      key: password
{{- end }}

{{- define "snooze.nats.env" }}
- name: NATS_URL
  value: "http://nats:4222"
{{- end }}

{{- define "snooze.rabbitmqEnv" }}
- name: AMQP_USERNAME
  valueFrom:
    secretKeyRef:
      name: "{{ .Release.Name }}-rabbitmq-default-user"
      key: username
- name: AMQP_PASSWORD
  valueFrom:
    secretKeyRef:
      name: "{{ .Release.Name }}-rabbitmq-default-user"
      key: password
{{- end }}

{{- define "snooze.image" }}
{{- if .Values.image.digest }}
{{ .Values.image.repo }}/snooze@{{ .Values.image.digest }}
{{- else }}
{{ .Values.image.repo }}/snooze:{{ .Values.image.tag | default .Chart.AppVersion }}
{{- end }}
{{- end }}
