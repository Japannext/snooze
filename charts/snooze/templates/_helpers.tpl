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
