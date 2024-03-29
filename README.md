# Small Golang CLI to generate README for Gitlab CI components

Very small CLI based on `spf13/cobra` to generate `README.md` from existing GitLab components.

The genreated `README.md` can be expanded using a header and footer file. By default those are 
`HEADER.md` and `FOOTER.md` in the project directory.

The same goes for each component, if the component is in it's own directory within `templates/`.
For components that only consist of a file (eg. `templates/component-name.yaml`), no header or footer
files will be used.

## Supported inputs
The following fields on each input are supported
- `description`
- `default`
- `options`
- `type`
- `regex`

## Example

The following `spec` in a component

```yaml
spec:
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
      default: true
```

Will result in the following markdown

```markdown
| Input / Variable | Description | Default value | Type    | Options | Regex |
| ---------------- | ----------- | ------------- | ------- | ------- | ----- |
| `concurrency`    |             | _1_           | number  | __ | `` |
| `environment`    |             | __            |         | _test, staging, production_ | `` |
| `export_results` |             | _true_        | boolean | __ | `` |
| `job-prefix`     | Define a prefix for the job name | __            |         | __ | `` |
| `job-stage`      |             | _test_        |         | __ | `` |
| `version`        |             | __            | string  | __ | `/^v\d\.\d+(\.\d+)$/` |
```

### Generated table
| Input / Variable | Description | Default value | Type    | Options | Regex |
| ---------------- | ----------- | ------------- | ------- | ------- | ----- |
| `concurrency`    |             | _1_           | number  | __ | `` |
| `environment`    |             | __            |         | _test, staging, production_ | `` |
| `export_results` |             | _true_        | boolean | __ | `` |
| `job-prefix`     | Define a prefix for the job name | __            |         | __ | `` |
| `job-stage`      |             | _test_        |         | __ | `` |
| `version`        |             | __            | string  | __ | `/^v\d\.\d+(\.\d+)$/` |
