package internal

import (
	"io"
	"text/template"
)

const enumerTmpl = fileHeaderText + `
package {{.PackageName}}

{{ range $enum := .Enums }}
// {{ $enum.EnumName }} enum declarations
{{ if eq $enum.TypeName "int" }}
{{- /* declarations specific for int valued enums */ -}}
const (
	_{{$enum.EnumName}}Name = "{{ range $item := $enum.Values }}{{ $item.Value }}{{ end }}"
{{- range $i, $item := $enum.Values }}
{{- if eq $i 0 }}
	{{ $item.Name }} {{ $enum.EnumName }} = iota
{{- else }}
	{{ $item.Name }}
{{- end }}
{{- end }}
)

var _{{ $enum.EnumName }}Names = []string{
{{- $prevIndex := 0 }}
{{- range $item := $enum.Values }}
	{{- $currIndex := sumWithLen $prevIndex $item.Value }}
	_{{ $enum.EnumName }}Name[{{ $prevIndex }}:{{ $currIndex }}],
	{{- $prevIndex = $currIndex }}
{{- end }}
}

var _{{ $enum.EnumName }}Map = map[string]{{ $enum.EnumName }}{
{{- $prevIndex := 0 }}
{{- range $item := $enum.Values }}
	{{- $currIndex := sumWithLen $prevIndex $item.Value }}
	_{{ $enum.EnumName }}Name[{{ $prevIndex }}:{{ $currIndex }}]: {{ $item.Name }},
	{{- $prevIndex = $currIndex }}
{{- end }}
}

func {{ $enum.EnumName }}Names() []string {
	return _{{ $enum.EnumName }}Names
}

func {{ $enum.EnumName }}FromString(value string) ({{ $enum.EnumName }}, bool) {
	enum, has := _{{ $enum.EnumName }}Map[value]
	return enum, has
}

func (en {{ $enum.EnumName }}) String() string {
	return _{{ $enum.EnumName }}Names[en]
}

{{- else if eq $enum.TypeName "string" -}}
{{- /* declarations specific for string valued enums */ -}}
const (
{{- range $i, $item := $enum.Values }}
	{{ $item.Name }} {{ $enum.EnumName }} = "{{ $item.Value }}"
{{- end }}
)

func (en {{ $enum.EnumName }}) String() string {
	return string(en)
}
{{- end }}

func {{ $enum.EnumName }}Values() []{{ $enum.EnumName }} {
	return []{{ $enum.EnumName }}{
	{{- range $item := $enum.Values }}
		{{ $item.Name }},
	{{- end }}
	}
}

func (en {{ $enum.EnumName }}) IsValid() bool {
	{{- range $i, $item := $enum.Values -}}
	{{- if eq $i 0 }}
	return en == {{ $item.Name }}{{ else }} ||
		en == {{ $item.Name }}
	{{- end }}
	{{- end }}
}

{{ end }}
`

func generateEnumerFileContent(packageName string, enums []*FutureEnum, w io.Writer) error {
	funcMap := template.FuncMap{
		"sumWithLen": func(one int, str string) int {
			return one + len(str)
		},
	}
	tmpl, err := template.New("enumer").Funcs(funcMap).Parse(enumerTmpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, map[string]any{
		"PackageName": packageName,
		"Enums":       enums,
	})
}
