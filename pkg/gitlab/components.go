package gitlab

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var readmeTemplate = `## {{ .Name }}

| Input / Variable | Description                            | Default value     | Options     |
| --------------------- | -------------------------------------- | ----------------- | ----------------- |
{{ range $key, $value := .Spec.Inputs -}}
	{{ $value.Markdown $key }}
{{ end }}

`

type ComponentInput struct {
	Default     string   `yaml:"default"`
	Description string   `yaml:"description"`
	Options     []string `yaml:"options"`
}

func (c ComponentInput) Markdown(name string) string {
	return fmt.Sprintf("| `%s` | %s | _%s_ | _%s_ |", name, c.Description, c.Default, strings.Join(c.Options, ", "))
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
