package gitlab

import (
	"bytes"
	"fmt"
	"text/template"
)

var readmeTemplate = `## {{ .Name }}

| Input / Variable | Description                            | Default value     |
| --------------------- | -------------------------------------- | ----------------- |
{{ range $key, $value := .Spec.Inputs -}}
	| ` + "`{{ $key }}`" + ` | {{ $value.Description }} | _{{ $value.Default }}_ |
{{ end }}

`

type ComponentInput struct {
	Default     string `yaml:"default"`
	Description string `yaml:"description"`
}

type ComponentSpec struct {
	Inputs map[string]ComponentInput `yaml:"inputs"`
}

type Component struct {
	Name string
	Spec ComponentSpec `yaml:"spec"`
}

func (c *Component) Markdown() string {

	// render go template
	t := template.Must(template.New("readme").Parse(readmeTemplate))
	var tpl bytes.Buffer
	err := t.Execute(&tpl, c)
	if err != nil {
		fmt.Printf("failed to render template: %v", err)
	}
	return tpl.String()
}
