{{- define "snooze.opensearch.env" }}
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

{{- define "snooze.proxy.env" }}
{{- if .Values.proxy }}
- name: HTTP_PROXY
  value: "{{ .Values.proxy }}"
- name: HTTPS_PROXY
  value: "{{ .Values.proxy }}"
- name: NO_PROXY
  value: "{{ .Values.noProxy }}"
{{- end }}
{{- end }}

{{- define "snooze.image" }}
{{- if .Values.image.digest }}
{{ .Values.image.repo }}/snooze@{{ .Values.image.digest }}
{{- else }}
{{ .Values.image.repo }}/snooze:{{ .Values.image.tag | default .Chart.AppVersion }}
{{- end }}
{{- end }}
