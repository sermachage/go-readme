package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunDoctor_ShowsActionableFailures(t *testing.T) {
	dir := t.TempDir()
	withWorkingDir(t, dir)

	c := &doctorCommandStub{}
	runDoctor(c.command(), nil)

	out := c.output()
	if !strings.Contains(out, "Some checks failed") {
		t.Fatalf("expected failed summary, got:\n%s", out)
	}
	if !strings.Contains(out, "missing go.mod") {
		t.Fatalf("expected go.mod hint, got:\n%s", out)
	}
	if !strings.Contains(out, "not a git repository") {
		t.Fatalf("expected git init hint, got:\n%s", out)
	}
	if !strings.Contains(out, "run `go-readme generate`") {
		t.Fatalf("expected README action hint, got:\n%s", out)
	}
}

func TestRunDoctor_AllChecksPass(t *testing.T) {
	dir := t.TempDir()
	withWorkingDir(t, dir)

	writeCmdFile(t, dir, "go.mod", "module github.com/example/project\n\ngo 1.24\n")
	writeCmdFile(t, dir, "README.md", "# Project\n")
	runGit(t, dir, "init")
	runGit(t, dir, "remote", "add", "origin", "https://github.com/example/project.git")

	c := &doctorCommandStub{}
	runDoctor(c.command(), nil)

	out := c.output()
	if !strings.Contains(out, "All checks passed.") {
		t.Fatalf("expected success summary, got:\n%s", out)
	}
	if !strings.Contains(out, "https://github.com/example/project.git") {
		t.Fatalf("expected remote URL detail, got:\n%s", out)
	}
}

type doctorCommandStub struct {
	out bytes.Buffer
	err bytes.Buffer
}

func (s *doctorCommandStub) command() *cobra.Command {
	c := &cobra.Command{}
	c.SetOut(&s.out)
	c.SetErr(&s.err)
	return c
}

func (s *doctorCommandStub) output() string {
	return s.out.String() + s.err.String()
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	c := exec.Command("git", args...)
	c.Dir = dir
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(out))
	}
}

func TestGitRepoStatus_WhenNotRepo(t *testing.T) {
	ok, detail := gitRepoStatus(t.TempDir())
	if ok {
		t.Fatal("expected gitRepoStatus to fail for non-repo dir")
	}
	if detail == "" {
		t.Fatal("expected detail for non-repo dir")
	}
}

func TestGitRemoteStatus_WhenNoRemote(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")

	ok, detail := gitRemoteStatus(dir)
	if ok {
		t.Fatal("expected gitRemoteStatus to fail when origin is missing")
	}
	if !strings.Contains(detail, "no remote.origin set") {
		t.Fatalf("unexpected detail: %q", detail)
	}
}

func TestGitRepoStatus_WhenRepo(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")

	ok, detail := gitRepoStatus(dir)
	if !ok {
		t.Fatalf("expected repo status true, detail=%q", detail)
	}
}

func TestGitRemoteStatus_WhenRemoteExists(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")
	runGit(t, dir, "remote", "add", "origin", "https://github.com/example/project.git")

	ok, detail := gitRemoteStatus(dir)
	if !ok {
		t.Fatalf("expected remote status true, detail=%q", detail)
	}
	if detail != "https://github.com/example/project.git" {
		t.Fatalf("unexpected remote detail: %q", detail)
	}
}

func TestGitOutput_Error(t *testing.T) {
	_, err := gitOutput(t.TempDir(), "not-a-real-git-subcommand")
	if err == nil {
		t.Fatal("expected gitOutput error")
	}
}

func TestGitOutput_Success(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")

	out, err := gitOutput(dir, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		t.Fatalf("gitOutput: %v", err)
	}
	if strings.TrimSpace(out) != "true" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestRunDoctor_DoesNotRequireReadmeDirectoryPath(t *testing.T) {
	dir := t.TempDir()
	withWorkingDir(t, dir)

	// Ensure we are validating the current working directory and not a hardcoded path.
	if _, err := os.Stat(filepath.Join(dir, "README.md")); !os.IsNotExist(err) {
		t.Fatalf("unexpected README in temp dir: %v", err)
	}

	c := &doctorCommandStub{}
	runDoctor(c.command(), nil)
	if !strings.Contains(c.output(), "README.md not found") {
		t.Fatalf("expected README not found message, got:\n%s", c.output())
	}
}
