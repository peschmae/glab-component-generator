/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peschmae/glab-component-generator/pkg/gitlab"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "readme",
		Aliases: []string{"r"},
		Short:   "Generates a README.md for all components within the given project directory",
		Long: `Gathers all components in <project>/templates and generates a README.md
from them using the inputs spec.

The generated README is prepended by a HEADER and FOOTER file, if present.
The same goes for each component.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateFlags()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateReadme()
		},
	}

	cmd.Flags().StringP("project", "p", ".", "The path to the gitlab CI component project")
	cmd.Flags().StringP("output", "o", "README.md", "The path to the output file. Relative to the projet directory")

	cmd.Flags().String("header", "HEADER.md", "File to prepended to the list of components")
	cmd.Flags().String("footer", "FOOTER.md", "File to appended to the list of components")

	cmd.Flags().String("component-header", "HEADER.md", "File to prepended on component. The file must exist in the component directory")
	cmd.Flags().String("component-footer", "FOOTER.md", "File to appended on component. The file must exist in the component directory")

	cmd.Flags().Int("component-header-level", 2, "The level of the header for each component")

	// bind flags to viper
	viper.BindPFlag("project", cmd.Flags().Lookup("project"))
	viper.BindPFlag("output", cmd.Flags().Lookup("output"))

	viper.BindPFlag("header", cmd.Flags().Lookup("header"))
	viper.BindPFlag("footer", cmd.Flags().Lookup("footer"))

	viper.BindPFlag("component-header", cmd.Flags().Lookup("component-header"))
	viper.BindPFlag("component-footer", cmd.Flags().Lookup("component-footer"))

	viper.BindPFlag("component-header-level", cmd.Flags().Lookup("component-header-level"))

	return cmd
}

func generateReadme() error {
	components := []string{}
	templatePath := filepath.Join(viper.GetString("project"), "templates")
	// find all yaml files in project
	filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// within the templates directory, we take all the yaml/yml files
		if filepath.Dir(path) == templatePath && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			components = append(components, path)
		} else if filepath.Dir(path) != templatePath && (filepath.Base(path) == "template.yaml" || filepath.Base(path) == "template.yml") {
			// if we are in a subdirectory, only the template.yaml/yml files are relevant
			components = append(components, path)
		}
		return nil
	})

	var sb strings.Builder
	if _, err := os.Stat(filepath.Join(viper.GetString("project"), viper.GetString("header"))); err == nil {
		header, err := os.ReadFile(filepath.Join(viper.GetString("project"), viper.GetString("header")))
		if err != nil {
			return err
		}
		sb.WriteString(string(header))
	} else {
		sb.WriteString("# GitLab CI Components\n\nThis repository contains the following components:\n\n[[_TOC_]]\n")
	}

	sb.WriteString("\n")
	// for each yaml file, parse and render the markdown
	for i := 0; i < len(components); i++ {
		c, err := gitlab.NewComponent(components[i])
		if err != nil {
			return err
		}
		// render markdown
		sb.WriteString(c.Markdown())
	}

	if _, err := os.Stat(filepath.Join(viper.GetString("project"), viper.GetString("footer"))); err == nil {
		footer, err := os.ReadFile(filepath.Join(viper.GetString("project"), viper.GetString("footer")))
		if err != nil {
			return err
		}
		sb.WriteString(string(footer))
	}

	sb.WriteString("\n")

	// write to file
	err := os.WriteFile(filepath.Join(viper.GetString("project"), viper.GetString("output")), []byte(strings.TrimSpace(sb.String())+"\n"), 0644)
	if err != nil {
		return err
	}
	return nil
}

func validateFlags() error {
	// Check if project exists
	if _, err := os.Stat(viper.GetString("project")); os.IsNotExist(err) {
		return fmt.Errorf("project does not exist")
	}

	return nil
}
