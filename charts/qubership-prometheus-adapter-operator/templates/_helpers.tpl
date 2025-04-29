{{/* vim: set filetype=mustache: */}}

{{/*
Expand the name of the chart.
*/}}
{{- define "qubership-prometheus-adapter-operator.name" -}}
  {{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "qubership-prometheus-adapter-operator.fullname" -}}
  {{- if .Values.fullnameOverride }}
    {{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
  {{- else }}
    {{- $name := default .Chart.Name .Values.nameOverride }}
    {{- if contains $name .Release.Name }}
      {{- .Release.Name | trunc 63 | trimSuffix "-" }}
    {{- else }}
      {{- printf "%s" $name | trunc 63 | trimSuffix "-" }}
    {{- end }}
  {{- end }}
{{- end }}

{{/*
Create common labels for each resource which is creating by this chart.
*/}}
{{- define "qubership-prometheus-adapter-operator.commonLabels" -}}
app.kubernetes.io/component: prometheus-adapter
app.kubernetes.io/part-of: monitoring
app.kubernetes.io/version: {{ .Chart.AppVersion }}
{{- end }}

{{/*
Find a prometheus-adapter-operator image in various places.
Image can be found from:
* specified by user from .Values.image
* default value
*/}}
{{- define "prometheus-adapter-operator.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- printf "ghcr.io/netcracker/qubership-prometheus-adapter-operator:main" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a prometheus-adapter image in various places.
Image can be found from:
* specified by user from .Values.prometheusAdapter.image
* default value
*/}}
{{- define "prometheus-adapter.image" -}}
  {{- if .Values.prometheusAdapter.image -}}
    {{- printf "%s" .Values.prometheusAdapter.image -}}
  {{- else -}}
    {{- printf "ghcr.io/netcracker/prometheus-adapter:main" -}}
  {{- end -}}
{{- end -}}