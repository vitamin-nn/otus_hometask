//nolint
package main

// Наверное, лучше было бы в отдельный файл положить шаблоны и использовать библиотеку,
// которая скомпилирует шаблоны в бинарный код. Но здесь такой задачи не было
func getTemplates() string {
	return `{{define "validation"}}
// Code generated automatically. DO NOT EDIT.
package {{.PackageName}}
import (
    "errors"
    "regexp"
)
var (
	ErrMin = errors.New("Value is less than min")
	ErrMax = errors.New("Value is greater than max")
	ErrLen = errors.New("Length validation error")
	ErrRegexp = errors.New("Regexp validation error")
	ErrEnum = errors.New("Enum validation error")
)
type ValidationError struct {
	Field string
	Err error
}
{{range $structName, $structFields := .Structures}}
    func (s {{$structName}}) Validate() (verr []ValidationError, err error) {
        {{range $field := $structFields}}
            {{$fieldVar := getFieldVar $field.Name }}
            {{if $field.Regexp.IsSet}}
                {{template "regexp_compile" $field}}
            {{end}}
            {{if $field.Type.IsSlice}}
            	for _, {{$fieldVar}} := range s.{{$field.Name}} {
            {{else}}
                {{$fieldVar}} := s.{{$field.Name}}
            {{end}}
            {{if $field.Regexp.IsSet}}
                {{template "regexp" $field}}
            {{end}}
            {{if $field.Min.IsSet}}
                {{template "min" $field}}
            {{end}}
            {{if $field.Max.IsSet}}
                {{template "max" $field}}
            {{end}}
            {{if $field.Len.IsSet}}
                {{template "len" $field}}
            {{end}}
            {{if $field.Enum.IsSet}}
                {{template "enum" $field}}
            {{end}}
            {{if $field.Type.IsSlice}}
            	}
            {{end}}
        {{end}}
        return
    }
{{end}}
{{end}}
{{define "min"}}
	if ({{getFieldVar .Name}} < {{.Min.Value}}) {
		ve := ValidationError{
			Field: "{{.Name}}",
			Err: ErrMin,
		}
		verr = append(verr, ve)
		{{if .Type.IsSlice}}
			break
		{{end}}
	}
{{end}}
{{define "max"}}
	if ({{getFieldVar .Name}} > {{.Max.Value}}) {
		ve := ValidationError{
			Field: "{{.Name}}",
			Err: ErrMax,
		}
		verr = append(verr, ve)
		{{if .Type.IsSlice}}
			break
		{{end}}
	}
{{end}}
{{define "len"}}
	if (len({{getFieldVar .Name}}) != {{.Len.Value}}) {
		ve := ValidationError{
			Field: "{{.Name}}",
			Err: ErrLen,
		}
		verr = append(verr, ve)
		{{if .Type.IsSlice}}
			break
		{{end}}
	}
{{end}}
{{define "enum"}}
	{{$name := getFieldVar .Name}}
	{{$list := .GetEnumList}}
	if({{range $i, $v := $list}}{{if gt $i 0}} && {{end}}{{$name}} != {{$v}}{{end}}) {
		ve := ValidationError{
			Field: "{{.Name}}",
			Err: ErrEnum,
		}
		verr = append(verr, ve)
		{{if .Type.IsSlice}}
			break
		{{end}}
	}
{{end}}
{{define "regexp_compile"}}
	var re = regexp.MustCompile(` + "`{{.Regexp.Value}}`" + `)
{{end}}
{{define "regexp"}}
	{{$name := getFieldVar .Name}}
	if (!re.MatchString({{$name}})) {
		ve := ValidationError{
			Field: "{{.Name}}",
			Err: ErrRegexp,
		}
		verr = append(verr, ve)
		{{if .Type.IsSlice}}
			break
		{{end}}
	}
{{end}}`
}
