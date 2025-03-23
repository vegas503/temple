{{- /*
    Example usage:
    $ cat template.tpl | USER=foo COLORS=red,green,blue FEATURES=one=y,two=n,three=y temple
*/ -}}

User name: {{ env "USER" }}

Colors:
{{- $colors := env "COLORS" | split "," }}
{{- range $c := $colors }}
  - {{ $c }}
{{- end }}

Features:
{{- $features := envdefault "FEATURES" "" | split "," }}
{{- range $feat := $features }}
    {{- $kv := (split "=" $feat) }}
    {{- $key := (index $kv 0) }}
    {{- $value := (index $kv 1) }}
    Feature "{{ $key }}" enabled: {{ if (istrue $value) }}YES{{ else }}NO{{ end }}
{{- end }}
