{{- /*
    Example usage:
    $ cat template.tpl | USER=foo COLORS=red,green,blue FEATURES=one=y,two=n,three=y temple
*/ -}}

User name: {{ env "USER" }}

Colors:
{{- $colors := (split (env "COLORS") ",") }}
{{- range $colors }}
  - {{ . }}
{{- end }}

Features:
{{- $features := (split (envdefault "FEATURES" "") ",") }}
{{- range $features }}
    {{- $kv := (split . "=") }}
    {{- $key := (index $kv 0) }}
    {{- $value := (index $kv 1) }}
    Feature "{{ $key }}" enabled: {{ if (istrue $value) }}YES{{ else }}NO{{ end }}
{{- end }}
