package gitlab

import (
	"fmt"
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
		expected.WriteString(fmt.Sprintf("| `environment`    |             | %c             |         | _test, staging, production_ | `` |\n", '\U000026D4'))
		expected.WriteString("| `export_results` |             | _true_        | boolean | __ | `` |\n")
		expected.WriteString(fmt.Sprintf("| `job-prefix`     | Define a prefix for the job name | %c             |         | __ | `` |\n", '\U000026D4'))
		expected.WriteString("| `job-stage`      |             | _test_        |         | __ | `` |\n")
		expected.WriteString(fmt.Sprintf("| `version`        |             | %c             | string  | __ | `/^v\\d\\.\\d+(\\.\\d+)$/` |\n", '\U000026D4'))

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
		expected.WriteString(fmt.Sprintf("| `job-prefix`     | Define a prefix for the job name | %c             |         | `` |\n", '\U000026D4'))
		expected.WriteString("| `job-stage`      |             | _test_        |         | `` |\n")
		expected.WriteString(fmt.Sprintf("| `version`        |             | %c             | string  | `/^v\\d\\.\\d+(\\.\\d+)$/` |\n", '\U000026D4'))

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
		expected.WriteString(fmt.Sprintf("| `environment`    |             | %c             |         | _test, staging, production_ |\n", '\U000026D4'))
		expected.WriteString("| `export_results` |             | _true_        | boolean | __ |\n")
		expected.WriteString(fmt.Sprintf("| `job-prefix`     | Define a prefix for the job name | %c             |         | __ |\n", '\U000026D4'))
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
		expected.WriteString(fmt.Sprintf("| `environment`    |             | %c             | _test, staging, production_ | `` |\n", '\U000026D4'))
		expected.WriteString(fmt.Sprintf("| `job-prefix`     | Define a prefix for the job name | %c             | __ | `` |\n", '\U000026D4'))
		expected.WriteString("| `job-stage`      |             | _test_        | __ | `` |\n")
		expected.WriteString(fmt.Sprintf("| `version`        |             | %c             | __ | `/^v\\d\\.\\d+(\\.\\d+)$/` |\n", '\U000026D4'))

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
		expected.WriteString(fmt.Sprintf("| `job-prefix`     | Define a prefix for the job name | %c             |\n", '\U000026D4'))
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
		expected.WriteString(fmt.Sprintf("| `job-prefix`     | Define a prefix for the job name | %c             |\n", '\U000026D4'))

		expected.WriteString("\n")

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
		expected.WriteString(fmt.Sprintf("| `job-prefix`     | Define a prefix for the job name | %c             |\n", '\U000026D4'))

		expected.WriteString("\n")

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
		expected.WriteString(fmt.Sprintf("| `job-prefix`     | Define a prefix for the job name | %c             |\n", '\U000026D4'))

		expected.WriteString("\n")

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
		expected.WriteString(fmt.Sprintf("| `job-prefix`     | Define a prefix for the job name | %c             |\n", '\U000026D4'))

		expected.WriteString("\nSome Footer\n")

		component := &Component{Name: "Footer test", Footer: "Some Footer"}
		yaml.Unmarshal([]byte(input), component)

		assert.Equal(t, expected.String(), component.Markdown())

	})

	t.Run("Header component level", func(t *testing.T) {

		viper.Set("component-header-level", 3)

		var expected strings.Builder
		expected.WriteString(`### Header level test

| Input / Variable | Description | Default value |
| ---------------- | ----------- | ------------- |
`)
		expected.WriteString(fmt.Sprintf("| `job-prefix`     | Define a prefix for the job name | %c             |\n", '\U000026D4'))

		expected.WriteString("\n")

		component := &Component{Name: "Header level test"}
		yaml.Unmarshal([]byte(input), component)

		assert.Equal(t, expected.String(), component.Markdown())

	})

}
