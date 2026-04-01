package template_test

import (
	"strings"
	"testing"

	"github.com/sermachage/go-readme/internal/domain"
	tmpl "github.com/sermachage/go-readme/internal/template"
)

func TestRender_GoDefault(t *testing.T) {
	renderer := tmpl.NewRenderer()
	project := domain.Project{
		Name:                   "myproject",
		ModulePath:             "github.com/example/myproject",
		GoVersion:              "1.21",
		RepoURL:                "https://github.com/example/myproject",
		Description:            "A great project",
		Features:               []string{"fast", "friendly"},
		UsageExample:           "import \"github.com/example/myproject\"",
		UsageLanguage:          "go",
		Configuration:          "Set MYPROJECT_ENV before running the CLI.",
		Dependencies:           []string{"github.com/spf13/cobra"},
		ContributingGuide:      "CONTRIBUTING.md",
		SecurityPolicy:         "SECURITY.md",
		License:                "LICENSE",
		AdditionalDependencies: 1,
	}

	got, err := renderer.Render("go_default.md", project)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	checks := []string{
		"# myproject",
		"A great project",
		"## Features",
		"friendly",
		"## Configuration",
		"MYPROJECT_ENV",
		"## Development",
		"github.com/example/myproject",
		"Go 1.21",
		"github.com/spf13/cobra",
		"https://github.com/example/myproject",
		"CONTRIBUTING.md",
		"SECURITY.md",
		"LICENSE",
	}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("rendered output missing %q\nfull output:\n%s", want, got)
		}
	}
}

func TestRender_UnknownTemplate(t *testing.T) {
	renderer := tmpl.NewRenderer()
	_, err := renderer.Render("nonexistent.md", domain.Project{})
	if err == nil {
		t.Fatal("expected error for nonexistent template, got nil")
	}
}
