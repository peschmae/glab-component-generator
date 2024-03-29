/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package gitlab

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var readmeTemplate = `## {{ .Name }}
{{ .Header }}
{{ if .HasOptions }}
| Input / Variable | Description | Default value | Options |
| ---------------- | ----------- | ------------- | ------- |
{{ range $key, $value := .Spec.Inputs -}}
	{{ $value.Markdown $key }}
{{ end }}
{{ else }}
| Input / Variable | Description | Default value |
| ---------------- | ----------- | ------------- |
{{ range $key, $value := .Spec.Inputs -}}
	{{ $value.MarkdownWithoutOptions $key }}
{{ end }}
{{ end }}
{{ .Footer }}
`

type ComponentInput struct {
	Default     string   `yaml:"default"`
	Description string   `yaml:"description"`
	Options     []string `yaml:"options"`
}

func (c ComponentInput) Markdown(name string) string {
	return fmt.Sprintf("| `%s` | %s | _%s_ | _%s_ |", name, c.Description, c.Default, strings.Join(c.Options, ", "))
}

func (c ComponentInput) MarkdownWithoutOptions(name string) string {
	return fmt.Sprintf("| `%s` | %s | _%s_ |", name, c.Description, c.Default)
}

type ComponentSpec struct {
	Inputs map[string]ComponentInput `yaml:"inputs"`
}

type Component struct {
	Name   string
	Header string
	Footer string
	Spec   ComponentSpec `yaml:"spec"`
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

func (c *Component) HasOptions() bool {

	for _, input := range c.Spec.Inputs {
		if len(input.Options) > 0 {
			return true
		}
	}
	return false
}

func NewComponent(path string) (*Component, error) {
	var name string
	var header []byte
	var footer []byte
	// GitLab allows yaml files directly in template directory, there we need to get the name from the filename
	// Otherwise the name is the parent directory name
	if filepath.Base(path) == "template.yml" || filepath.Base(path) == "template.yaml" {
		name = filepath.Base(filepath.Dir(path))

		if _, err := os.Stat(filepath.Join(filepath.Dir(path), viper.GetString("component-header"))); err == nil {
			header, err = os.ReadFile(filepath.Join(filepath.Dir(path), viper.GetString("component-header")))
			if err != nil {
				return nil, err
			}
		}

		if _, err := os.Stat(filepath.Join(filepath.Dir(path), viper.GetString("component-footer"))); err == nil {
			footer, err = os.ReadFile(filepath.Join(filepath.Dir(path), viper.GetString("component-footer")))
			if err != nil {
				return nil, err
			}
		}
	} else {
		name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	c := &Component{Name: name, Header: string(header), Footer: string(footer)}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	yaml.Unmarshal(b, c)

	return c, nil
}
