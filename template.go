package main

import "text/template"

var generatedTmpl = template.Must(template.New("generated").Parse(`
Todo's
======

{{ $length := len .Todos }} {{ if eq $length 0 }}
### Hurray! Nothing to do
{{end}}

{{range $index, $todo := .Todos}}
- [ ] {{$todo.Text}} ([{{$todo.Pos.Filename}}]({{$todo.Pos.Filename}}) {{$todo.Pos.Line}}:{{$todo.Pos.Column}})
{{end}}

> Generated with todos {{.Command}}, for more information: [todos](https://github.com/onethousandone/todos)
`))
