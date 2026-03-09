package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunGenerate_DryRunWritesToCommandOutput(t *testing.T) {
	dir := t.TempDir()
	writeCmdFile(t, dir, "go.mod", "module github.com/example/dryrun\n\ngo 1.24\n")
	resetGenerateFlags(t)
	generateDir = dir

	generateDescription = "dry run project"
	generateDryRun = true

	c := &cobraCommandStub{}
	if err := runGenerate(c.command(), nil); err != nil {
		t.Fatalf("runGenerate: %v", err)
	}

	out := c.output()
	if !strings.Contains(out, "# dryrun") {
		t.Fatalf("expected generated README in command output, got:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(dir, "README.md")); !os.IsNotExist(err) {
		t.Fatalf("README.md exists after dry run, err=%v", err)
	}
}

func TestRunGenerate_CreateReadme(t *testing.T) {
	dir := t.TempDir()
	writeCmdFile(t, dir, "go.mod", "module github.com/example/project\n\ngo 1.24\n")
	resetGenerateFlags(t)
	generateDir = dir

	generateDescription = "new project"

	c := &cobraCommandStub{}
	if err := runGenerate(c.command(), nil); err != nil {
		t.Fatalf("runGenerate: %v", err)
	}

	out := c.output()
	if !strings.Contains(out, "README created:") {
		t.Fatalf("expected create message, got: %s", out)
	}
	got, err := os.ReadFile(filepath.Join(dir, "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	if !strings.Contains(string(got), "new project") {
		t.Fatalf("README missing expected description, got:\n%s", string(got))
	}
}

func TestRunGenerate_UsesConfiguredDir(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "project")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeCmdFile(t, target, "go.mod", "module github.com/example/from-flag\n\ngo 1.24\n")
	resetGenerateFlags(t)
	generateDir = target
	generateDescription = "from dir flag"
	generateNonInteractive = true

	c := &cobraCommandStub{}
	if err := runGenerate(c.command(), nil); err != nil {
		t.Fatalf("runGenerate: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(target, "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	if !strings.Contains(string(got), "from dir flag") {
		t.Fatalf("README missing expected description, got:\n%s", string(got))
	}
}

func TestPromptDescription_NonTTYReturnsEmpty(t *testing.T) {
	originalStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = originalStdin
		_ = r.Close()
		_ = w.Close()
	})

	c := &cobraCommandStub{}
	got := promptDescription(c.command())
	if got != "" {
		t.Fatalf("promptDescription = %q, want empty string for non-tty stdin", got)
	}
	if c.output() != "" {
		t.Fatalf("expected no prompt output for non-tty stdin, got: %q", c.output())
	}
}

type cobraCommandStub struct {
	out bytes.Buffer
}

func (s *cobraCommandStub) command() *cobra.Command {
	c := &cobra.Command{}
	c.SetOut(&s.out)
	return c
}

func (s *cobraCommandStub) output() string {
	return s.out.String()
}

func writeCmdFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(prev)
	})
}

func resetGenerateFlags(t *testing.T) {
	t.Helper()
	generateDescription = ""
	generateDir = "."
	generateTemplate = "go_default.md"
	generateDryRun = false
	generateForce = false
	generateNonInteractive = false
	t.Cleanup(func() {
		generateDescription = ""
		generateDir = "."
		generateTemplate = "go_default.md"
		generateDryRun = false
		generateForce = false
		generateNonInteractive = false
	})
}
