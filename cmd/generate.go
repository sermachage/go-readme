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
	generateDir            string
	generateDescription    string
	generateFeatures       []string
	generateUsageExample   string
	generateConfiguration  string
	generateContributing   string
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
	generateCmd.Flags().StringVar(&generateDir, "dir", ".", "target project directory")
	generateCmd.Flags().StringVarP(&generateDescription, "description", "d", "", "project description")
	generateCmd.Flags().StringSliceVar(&generateFeatures, "features", nil, "comma-separated key project features")
	generateCmd.Flags().StringVar(&generateUsageExample, "usage-example", "", "usage example command or snippet")
	generateCmd.Flags().StringVar(&generateConfiguration, "configuration", "", "configuration notes")
	generateCmd.Flags().StringVar(&generateContributing, "contributing-notes", "", "contributing guidance")
	generateCmd.Flags().StringVarP(&generateTemplate, "template", "t", "go_default.md", "template file name")
	generateCmd.Flags().BoolVar(&generateDryRun, "dry-run", false, "print the README without writing to disk")
	generateCmd.Flags().BoolVar(&generateForce, "force", false, "overwrite the entire README (ignore markers)")
	generateCmd.Flags().BoolVar(&generateNonInteractive, "non-interactive", false, "disable interactive prompts")
}

func runGenerate(cmd *cobra.Command, _ []string) error {
	answers := generatePromptAnswers{
		Description:   strings.TrimSpace(generateDescription),
		Features:      append([]string(nil), generateFeatures...),
		UsageExample:  strings.TrimSpace(generateUsageExample),
		Configuration: strings.TrimSpace(generateConfiguration),
		Contributing:  strings.TrimSpace(generateContributing),
	}
	if !generateNonInteractive {
		answers = promptGenerateMetadata(cmd, answers)
	}

	opts := app.GenerateOptions{
		Dir:           generateDir,
		Description:   answers.Description,
		Features:      answers.Features,
		UsageExample:  answers.UsageExample,
		Configuration: answers.Configuration,
		Contributing:  answers.Contributing,
		Template:      generateTemplate,
		DryRun:        generateDryRun,
		Force:         generateForce,
	}

	result, err := app.Generate(opts)
	if err != nil {
		return err
	}

	if generateDryRun {
		fmt.Fprintln(cmd.OutOrStdout(), result.Content)
		return nil
	}

	action := "updated"
	if result.Created {
		action = "created"
	}
	fmt.Fprintf(cmd.OutOrStdout(), "README %s: %s\n", action, result.OutputPath)
	return nil
}

type generatePromptAnswers struct {
	Description   string
	Features      []string
	UsageExample  string
	Configuration string
	Contributing  string
}

// promptGenerateMetadata asks the user for optional README metadata on stdin.
// Returns the existing answers unchanged if stdin is not a terminal.
func promptGenerateMetadata(cmd *cobra.Command, answers generatePromptAnswers) generatePromptAnswers {
	if !stdinIsInteractive() {
		return answers
	}

	scanner := bufio.NewScanner(os.Stdin)
	if answers.Description == "" {
		answers.Description = promptLine(cmd, scanner, "Project description (leave blank to auto-detect/skip): ")
	}
	if len(answers.Features) == 0 {
		answers.Features = parsePromptList(promptLine(cmd, scanner, "Key features (comma-separated, leave blank to skip): "))
	}
	if answers.UsageExample == "" {
		answers.UsageExample = promptLine(cmd, scanner, "Usage example or command (leave blank for default): ")
	}
	if answers.Configuration == "" {
		answers.Configuration = promptLine(cmd, scanner, "Configuration notes (leave blank to skip): ")
	}
	if answers.Contributing == "" {
		answers.Contributing = promptLine(cmd, scanner, "Contributing notes (leave blank to skip): ")
	}
	return answers
}

func stdinIsInteractive() bool {
	fi, err := os.Stdin.Stat()
	return err == nil && (fi.Mode()&os.ModeCharDevice) != 0
}

func promptLine(cmd *cobra.Command, scanner *bufio.Scanner, prompt string) string {
	fmt.Fprint(cmd.OutOrStdout(), prompt)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

func parsePromptList(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '\n'
	})
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil
	}
	return items
}
