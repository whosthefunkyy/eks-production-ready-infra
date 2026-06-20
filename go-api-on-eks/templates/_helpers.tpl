{{- define "go-api.name" -}}
go-api
{{- end }}

{{- define "go-api.labels" -}}
app: {{ include "go-api.name" . }}
project: my-k8s-infra
{{- end }}
