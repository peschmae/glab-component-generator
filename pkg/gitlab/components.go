/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package gitlab

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type ComponentInput struct {
	Default     string   `yaml:"default"`
	Description string   `yaml:"description"`
	Options     []string `yaml:"options"`
	Type        string   `yaml:"type"`
	Regex       string   `yaml:"regex"`
}

func headerLevel() string {
	return strings.Repeat("#", viper.GetInt("component-header-level"))
}

// replace linebreaks with <br> as Gitlab converts it to HTML anyways
func replaceLinebreaks(input string) string {
	return strings.TrimSuffix(strings.ReplaceAll(input, "\n", "<br>"), "<br>")
}

func (input ComponentInput) Markdown(name string, hasTypes, hasOptions, hasRegex bool) string {
	var sb strings.Builder
	if input.Default == "" {
		input.Default = fmt.Sprintf("%c", '\U000026D4')
	} else {
		input.Default = fmt.Sprintf("_%s_", input.Default)
	}
	sb.WriteString(fmt.Sprintf("| %-16s | %-11s | %-13s |", fmt.Sprintf("`%s`", name), replaceLinebreaks(input.Description), input.Default))

	if hasTypes {
		sb.WriteString(fmt.Sprintf(" %-7s |", input.Type))
	}
	if hasOptions {
		sb.WriteString(fmt.Sprintf(" _%s_ |", strings.Join(input.Options, ", ")))
	}
	if hasRegex {
		sb.WriteString(fmt.Sprintf(" `%s` |", input.Regex))
	}
	sb.WriteString("\n")

	return sb.String()

}

type ComponentSpec struct {
	Inputs map[string]ComponentInput `yaml:"inputs"`
}

func (spec *ComponentSpec) MarkdownTable() string {
	hasTypes := spec.HasTypes()
	hasOptions := spec.HasOptions()
	hasRegex := spec.HasRegex()

	var sb strings.Builder
	var divider strings.Builder

	// Generate header
	sb.WriteString("| Input / Variable | Description | Default value |")
	divider.WriteString("| ---------------- | ----------- | ------------- |")
	if hasTypes {
		sb.WriteString(" Type    |")
		divider.WriteString(" ------- |")
	}
	if hasOptions {
		sb.WriteString(" Options |")
		divider.WriteString(" ------- |")
	}
	if hasRegex {
		sb.WriteString(" Regex |")
		divider.WriteString(" ----- |")
	}
	sb.WriteString("\n")
	// Write divider
	sb.WriteString(divider.String())
	sb.WriteString("\n")

	keys := make([]string, len(spec.Inputs))

	i := 0
	for k := range spec.Inputs {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for i = 0; i < len(keys); i++ {
		sb.WriteString(spec.Inputs[keys[i]].Markdown(keys[i], hasTypes, hasOptions, hasRegex))
	}

	return sb.String()

}

func (spec *ComponentSpec) HasOptions() bool {

	for _, input := range spec.Inputs {
		if len(input.Options) > 0 {
			return true
		}
	}
	return false
}

func (spec *ComponentSpec) HasTypes() bool {

	for _, input := range spec.Inputs {
		if input.Type != "" {
			return true
		}
	}
	return false
}

func (spec *ComponentSpec) HasRegex() bool {

	for _, input := range spec.Inputs {
		if input.Regex != "" {
			return true
		}
	}
	return false
}

type Component struct {
	Name   string
	Header string
	Footer string
	Spec   *ComponentSpec `yaml:"spec"`
}

func (c *Component) Markdown() string {

	if c.Header == "" && c.Footer == "" && c.Spec == nil {
		return ""
	}

	var md strings.Builder

	// render component header
	md.WriteString(fmt.Sprintf("%s %s\n\n", headerLevel(), c.Name))

	if c.Header != "" {
		md.WriteString(strings.TrimSpace(c.Header) + "\n\n")
	}

	if c.Spec != nil {
		md.WriteString(c.Spec.MarkdownTable() + "\n")
	}

	if c.Footer != "" {
		md.WriteString(strings.TrimSpace(c.Footer) + "\n")
	}

	return md.String()
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
