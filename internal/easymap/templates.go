package easymap

const mappingTemplate = `package main
	{{- $nameIn:=.NameIn }}
	{{- $structName:=.NameIn }}
	func MapTo{{ .StructType }} (in *{{ .InStructType }}) *{{ .StructType }} {
		out := &{{ .StructType }}{
			{{ range $field := .ListSimpleFields -}}
				{{- $field.Name }}: in.{{ $field.Name }},
			{{ end -}}
		}
		{{- .Content }}
		return out
	}`

const ifConditionPointerTemplate = `{{- $nameIn:=.NameIn }}
	{{- $structName:=.NameIn }}
	{{- $parentInName:=.ParentInName }}
	{{- $parentOutName:=.ParentOutName }}
	if {{$parentInName}}{{$nameIn}} != nil {
		{{$parentOutName}}{{$nameIn}} = &{{.StructType}}{
			{{ range $field := .ListSimpleFields }}
				{{- $field.Name}}: {{$parentInName}}{{$nameIn}}.{{$field.Name}},
			{{ end -}}
		}
		{{- .Content }}
	}`

const slicePointerTemplate = `{{- $nameIn:=.NameIn }}
	{{- $structName:=.NameIn }}
	{{- $parentInName:=.ParentInName }}
	{{- $parentOutName:=.ParentOutName }}
	var new{{$nameIn}}Slice []*{{.StructType}}
	for _, {{$nameIn}}Item := range {{$parentInName}}{{$nameIn}} {
		new{{$nameIn}} := &{{.StructType}}{
			{{range $field := .ListSimpleFields }}
				{{- $field.Name}}: {{$nameIn}}Item.{{$field.Name}},
			{{ end -}}
		}
		{{- .Content  }}
		new{{$nameIn}}Slice = append(new{{$nameIn}}Slice, new{{$nameIn}})
	}
	{{$parentOutName}}{{$structName}} = new{{$nameIn}}Slice`
