package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose common project configuration issues",
	Long: `doctor checks whether the current directory is correctly set up for
readmeaker: go.mod present, git initialised, remote configured, and README valid.`,
	Run: runDoctor,
}

func runDoctor(cmd *cobra.Command, _ []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), "error: cannot determine working directory:", err)
		return
	}

	allOK := true
	check := func(label, detail string, ok bool) {
		status := "✓"
		if !ok {
			status = "✗"
			allOK = false
		}
		fmt.Fprintf(cmd.OutOrStdout(), "  %s  %s", status, label)
		if detail != "" {
			fmt.Fprintf(cmd.OutOrStdout(), " (%s)", detail)
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	fmt.Fprintln(cmd.OutOrStdout(), "readmeaker doctor")
	fmt.Fprintln(cmd.OutOrStdout())

	// go.mod
	_, gomodErr := os.Stat(filepath.Join(dir, "go.mod"))
	check("go.mod present", "", gomodErr == nil)

	// git initialised
	_, gitErr := os.Stat(filepath.Join(dir, ".git"))
	check("git initialized", "", gitErr == nil)

	// remote configured
	check("git remote.origin configured", "", hasGitRemote(dir))

	// README
	_, readmeErr := os.Stat(filepath.Join(dir, "README.md"))
	check("README.md exists", "", readmeErr == nil)

	fmt.Fprintln(cmd.OutOrStdout())
	if allOK {
		fmt.Fprintln(cmd.OutOrStdout(), "All checks passed.")
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), "Some checks failed – see above.")
	}
}

func hasGitRemote(dir string) bool {
	c := exec.Command("git", "config", "--get", "remote.origin.url")
	c.Dir = dir
	return c.Run() == nil
}
