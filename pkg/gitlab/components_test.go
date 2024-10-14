package gitlab

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_ComponentSpecMarkdown(t *testing.T) {
	t.Run("All fields", func(t *testing.T) {
		input := `
inputs:
  job-prefix:     # Mandatory string input
    description: "Define a prefix for the job name"
  job-stage:      # Optional string input with a default value when not provided
    default: test
  environment:    # Mandatory input that must match one of the options
    options: ['test', 'staging', 'production']
  concurrency:
    type: number  # Optional numeric input with a default value when not provided
    default: 1
  version:        # Mandatory string input that must match the regular expression
    type: string
    regex: /^v\d\.\d+(\.\d+)$/
  export_results: # Optional boolean input with a default value when not provided
    type: boolean
    default: true`
		var expected strings.Builder
		expected.WriteString(`| Input / Variable | Description | Default value | Type    | Options | Regex |
| ---------------- | ----------- | ------------- | ------- | ------- | ----- |
`)
		expected.WriteString("| `concurrency`    |             | _1_           | number  | __ | `` |\n")
		expected.WriteString("| `environment`    |             | __            |         | _test, staging, production_ | `` |\n")
		expected.WriteString("| `export_results` |             | _true_        | boolean | __ | `` |\n")
		expected.WriteString("| `job-prefix`     | Define a prefix for the job name | __            |         | __ | `` |\n")
		expected.WriteString("| `job-stage`      |             | _test_        |         | __ | `` |\n")
		expected.WriteString("| `version`        |             | __            | string  | __ | `/^v\\d\\.\\d+(\\.\\d+)$/` |\n")

		spec := &ComponentSpec{}
		yaml.Unmarshal([]byte(input), spec)

		assert.Equal(t, expected.String(), spec.MarkdownTable())

	})

	t.Run("Without options", func(t *testing.T) {
		input := `
inputs:
  job-prefix:     # Mandatory string input
    description: "Define a prefix for the job name"
  job-stage:      # Optional string input with a default value when not provided
    default: test
  concurrency:
    type: number  # Optional numeric input with a default value when not provided
    default: 1
  version:        # Mandatory string input that must match the regular expression
    type: string
    regex: /^v\d\.\d+(\.\d+)$/
  export_results: # Optional boolean input with a default value when not provided
    type: boolean
    default: true`
		var expected strings.Builder
		expected.WriteString(`| Input / Variable | Description | Default value | Type    | Regex |
| ---------------- | ----------- | ------------- | ------- | ----- |
`)
		expected.WriteString("| `concurrency`    |             | _1_           | number  | `` |\n")
		expected.WriteString("| `export_results` |             | _true_        | boolean | `` |\n")
		expected.WriteString("| `job-prefix`     | Define a prefix for the job name | __            |         | `` |\n")
		expected.WriteString("| `job-stage`      |             | _test_        |         | `` |\n")
		expected.WriteString("| `version`        |             | __            | string  | `/^v\\d\\.\\d+(\\.\\d+)$/` |\n")

		spec := &ComponentSpec{}
		yaml.Unmarshal([]byte(input), spec)

		assert.Equal(t, expected.String(), spec.MarkdownTable())

	})

	t.Run("Without regex", func(t *testing.T) {
		input := `
inputs:
  job-prefix:     # Mandatory string input
    description: "Define a prefix for the job name"
  job-stage:      # Optional string input with a default value when not provided
    default: test
  environment:    # Mandatory input that must match one of the options
    options: ['test', 'staging', 'production']
  concurrency:
    type: number  # Optional numeric input with a default value when not provided
    default: 1
  export_results: # Optional boolean input with a default value when not provided
    type: boolean
    default: true`
		var expected strings.Builder
		expected.WriteString(`| Input / Variable | Description | Default value | Type    | Options |
| ---------------- | ----------- | ------------- | ------- | ------- |
`)
		expected.WriteString("| `concurrency`    |             | _1_           | number  | __ |\n")
		expected.WriteString("| `environment`    |             | __            |         | _test, staging, production_ |\n")
		expected.WriteString("| `export_results` |             | _true_        | boolean | __ |\n")
		expected.WriteString("| `job-prefix`     | Define a prefix for the job name | __            |         | __ |\n")
		expected.WriteString("| `job-stage`      |             | _test_        |         | __ |\n")

		spec := &ComponentSpec{}
		yaml.Unmarshal([]byte(input), spec)

		assert.Equal(t, expected.String(), spec.MarkdownTable())

	})

	t.Run("Without type", func(t *testing.T) {
		input := `
inputs:
  job-prefix:     # Mandatory string input
    description: "Define a prefix for the job name"
  job-stage:      # Optional string input with a default value when not provided
    default: test
  environment:    # Mandatory input that must match one of the options
    options: ['test', 'staging', 'production']
  version:        # Mandatory string input that must match the regular expression
    regex: /^v\d\.\d+(\.\d+)$/`
		var expected strings.Builder
		expected.WriteString(`| Input / Variable | Description | Default value | Options | Regex |
| ---------------- | ----------- | ------------- | ------- | ----- |
`)
		expected.WriteString("| `environment`    |             | __            | _test, staging, production_ | `` |\n")
		expected.WriteString("| `job-prefix`     | Define a prefix for the job name | __            | __ | `` |\n")
		expected.WriteString("| `job-stage`      |             | _test_        | __ | `` |\n")
		expected.WriteString("| `version`        |             | __            | __ | `/^v\\d\\.\\d+(\\.\\d+)$/` |\n")

		spec := &ComponentSpec{}
		yaml.Unmarshal([]byte(input), spec)

		assert.Equal(t, expected.String(), spec.MarkdownTable())

	})

	t.Run("Minimal", func(t *testing.T) {
		input := `
inputs:
  job-prefix:     # Mandatory string input
    description: "Define a prefix for the job name"
  job-stage:      # Optional string input with a default value when not provided
    default: test
  concurrency:
    default: 1
  export_results: # Optional boolean input with a default value when not provided
    default: true`
		var expected strings.Builder
		expected.WriteString(`| Input / Variable | Description | Default value |
| ---------------- | ----------- | ------------- |
`)
		expected.WriteString("| `concurrency`    |             | _1_           |\n")
		expected.WriteString("| `export_results` |             | _true_        |\n")
		expected.WriteString("| `job-prefix`     | Define a prefix for the job name | __            |\n")
		expected.WriteString("| `job-stage`      |             | _test_        |\n")

		spec := &ComponentSpec{}
		yaml.Unmarshal([]byte(input), spec)

		assert.Equal(t, expected.String(), spec.MarkdownTable())

	})

	t.Run("Description linebreak", func(t *testing.T) {
		input := `
inputs:
  export_results: # Optional boolean input with a default value when not provided
    default: true
    description: |
      Line 1
      Line 2`
		var expected strings.Builder
		expected.WriteString(`| Input / Variable | Description | Default value |
| ---------------- | ----------- | ------------- |
`)
		expected.WriteString("| `export_results` | Line 1<br>Line 2 | _true_        |\n")

		spec := &ComponentSpec{}
		yaml.Unmarshal([]byte(input), spec)

		assert.Equal(t, expected.String(), spec.MarkdownTable())

	})
}

func Test_ComponentMarkdown(t *testing.T) {
	viper.Set("component-header-level", 2)

	input := `
spec:
  inputs:
    job-prefix:
      description: "Define a prefix for the job name"`

	t.Run("Minimal", func(t *testing.T) {

		var expected strings.Builder
		expected.WriteString(`## Component test



| Input / Variable | Description | Default value |
| ---------------- | ----------- | ------------- |
`)
		expected.WriteString("| `job-prefix`     | Define a prefix for the job name | __            |\n")

		expected.WriteString("\n\n")

		component := &Component{Name: "Component test"}
		yaml.Unmarshal([]byte(input), component)

		assert.Equal(t, expected.String(), component.Markdown())

	})

	t.Run("Header", func(t *testing.T) {

		var expected strings.Builder
		expected.WriteString(`## Header test

Some Header

| Input / Variable | Description | Default value |
| ---------------- | ----------- | ------------- |
`)
		expected.WriteString("| `job-prefix`     | Define a prefix for the job name | __            |\n")

		expected.WriteString("\n\n")

		component := &Component{Name: "Header test", Header: "Some Header"}
		yaml.Unmarshal([]byte(input), component)

		assert.Equal(t, expected.String(), component.Markdown())

	})

	t.Run("Header with linebreak", func(t *testing.T) {

		var expected strings.Builder
		expected.WriteString(`## Header test

Some
Header

| Input / Variable | Description | Default value |
| ---------------- | ----------- | ------------- |
`)
		expected.WriteString("| `job-prefix`     | Define a prefix for the job name | __            |\n")

		expected.WriteString("\n\n")

		component := &Component{Name: "Header test", Header: "Some\nHeader\n"}
		yaml.Unmarshal([]byte(input), component)

		assert.Equal(t, expected.String(), component.Markdown())

	})

	t.Run("Footer", func(t *testing.T) {

		var expected strings.Builder
		expected.WriteString(`## Footer test

| Input / Variable | Description | Default value |
| ---------------- | ----------- | ------------- |
`)
		expected.WriteString("| `job-prefix`     | Define a prefix for the job name | __            |\n")

		expected.WriteString("\nSome Footer")

		component := &Component{Name: "Footer test", Footer: "Some Footer"}
		yaml.Unmarshal([]byte(input), component)

		assert.Equal(t, expected.String(), component.Markdown())

	})

}
