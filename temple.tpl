{{- $colors := (split (env "COLORS") ",") }}
{{- $numbers := (split (env "NUMBERS") ",") }}
{{- range (append $colors $numbers) }}
  * {{ . }}
{{- end }}
