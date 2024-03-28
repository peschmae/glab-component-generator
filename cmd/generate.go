/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
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

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a README.md for the specified GitLab component project",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateReadme()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("project", "p", ".", "The path to the gitlab CI component project")
	generateCmd.Flags().StringP("output", "o", "README.md", "The path to the output file")

	// bind flags to viper
	viper.BindPFlag("project", generateCmd.Flags().Lookup("project"))
	viper.BindPFlag("output", generateCmd.Flags().Lookup("output"))
}

func generateReadme() error {
	components := []string{}
	// find all yaml files in project
	filepath.Walk(filepath.Join(viper.GetString("project"), "templates"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			components = append(components, path)
		}
		return nil
	})

	var sb strings.Builder
	if _, err := os.Stat(filepath.Join(viper.GetString("project"), "HEADER.md")); err == nil {
		header, err := os.ReadFile(filepath.Join(viper.GetString("project"), "HEADER.md"))
		if err != nil {
			return err
		}
		sb.WriteString(string(header))
	} else {
		sb.WriteString("# GitLab CI Components\n\nThis repository contains the following components:\n\n[[_TOC_]]")
	}

	sb.WriteString("\n\n")
	// for each yaml file, parse and render the markdown
	for i := 0; i < len(components); i++ {
		c, err := gitlab.NewComponent(components[i])
		if err != nil {
			return err
		}
		// render markdown
		sb.WriteString(c.Markdown())
	}

	if _, err := os.Stat(filepath.Join(viper.GetString("project"), "FOOTER.md")); err == nil {
		footer, err := os.ReadFile(filepath.Join(viper.GetString("project"), "FOOTER.md"))
		if err != nil {
			return err
		}
		sb.WriteString(string(footer))
	}

	sb.WriteString("\n")

	// write to file
	err := os.WriteFile(filepath.Join(viper.GetString("project"), viper.GetString("output")), []byte(strings.TrimSpace(sb.String())), 0644)
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
