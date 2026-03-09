package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sermachage/go-readme/internal/gitmeta"
)

var doctorDir string

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose common project configuration issues",
	Long: `doctor checks whether the current directory is correctly set up for
go-readme: go.mod present, git initialised, remote configured, and README valid.`,
	Run: runDoctor,
}

func init() {
	doctorCmd.Flags().StringVar(&doctorDir, "dir", ".", "target project directory")
}

func runDoctor(cmd *cobra.Command, _ []string) {
	dir := doctorDir

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

	fmt.Fprintln(cmd.OutOrStdout(), "go-readme doctor")
	fmt.Fprintln(cmd.OutOrStdout())

	// go.mod
	gomodPath := filepath.Join(dir, "go.mod")
	_, gomodErr := os.Stat(gomodPath)
	gomodDetail := ""
	if gomodErr != nil {
		gomodDetail = "missing go.mod at " + gomodPath
	}
	check("go.mod present", gomodDetail, gomodErr == nil)

	// git initialized/worktree
	isRepo, repoDetail := gitRepoStatus(dir)
	check("git initialized", repoDetail, isRepo)

	// remote configured
	hasRemote, remoteDetail := gitRemoteStatus(dir)
	check("git remote.origin configured", remoteDetail, hasRemote)

	// README
	readmePath := filepath.Join(dir, "README.md")
	_, readmeErr := os.Stat(readmePath)
	readmeDetail := ""
	if readmeErr != nil {
		readmeDetail = "README.md not found; run `go-readme generate`"
	}
	check("README.md exists", readmeDetail, readmeErr == nil)

	fmt.Fprintln(cmd.OutOrStdout())
	if allOK {
		fmt.Fprintln(cmd.OutOrStdout(), "All checks passed.")
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), "Some checks failed – see above.")
	}
}

func gitRepoStatus(dir string) (bool, string) {
	isRepo, err := gitmeta.IsRepository(dir)
	if err != nil {
		return false, "not a git repository (`git init` to create one)"
	}
	if !isRepo {
		return false, "directory is outside a git work tree"
	}
	return true, ""
}

func gitRemoteStatus(dir string) (bool, string) {
	remote, err := gitmeta.RemoteOriginURL(dir)
	if err != nil {
		return false, "no remote.origin set (`git remote add origin <url>`)"
	}
	if remote == "" {
		return false, "remote.origin is empty"
	}
	return true, remote
}
