package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sermachage/go-readme/internal/app"
)

var (
	generateDescription    string
	generateTemplate       string
	generateDryRun         bool
	generateForce          bool
	generateNonInteractive bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate or update README.md for the current Go project",
	Long: `generate detects the Go project in the current directory, extracts
metadata from go.mod and git, renders a README template, and writes (or
idempotently updates) README.md.`,
	RunE: runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&generateDescription, "description", "d", "", "project description")
	generateCmd.Flags().StringVarP(&generateTemplate, "template", "t", "go_default.md", "template file name")
	generateCmd.Flags().BoolVar(&generateDryRun, "dry-run", false, "print the README without writing to disk")
	generateCmd.Flags().BoolVar(&generateForce, "force", false, "overwrite the entire README (ignore markers)")
	generateCmd.Flags().BoolVar(&generateNonInteractive, "non-interactive", false, "disable interactive prompts")
}

func runGenerate(cmd *cobra.Command, _ []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	description := generateDescription
	if description == "" && !generateNonInteractive {
		description = promptDescription(cmd)
	}

	opts := app.GenerateOptions{
		Dir:         dir,
		Description: description,
		Template:    generateTemplate,
		DryRun:      generateDryRun,
		Force:       generateForce,
	}

	result, err := app.Generate(opts)
	if err != nil {
		return err
	}

	if generateDryRun {
		fmt.Println(result.Content)
		return nil
	}

	action := "updated"
	if result.Created {
		action = "created"
	}
	fmt.Fprintf(cmd.OutOrStdout(), "README %s: %s\n", action, result.OutputPath)
	return nil
}

// promptDescription asks the user for a short project description on stdin.
// Returns an empty string if stdin is not a terminal or the user skips.
func promptDescription(cmd *cobra.Command) string {
	fmt.Fprint(cmd.OutOrStdout(), "Project description (leave blank to skip): ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}
