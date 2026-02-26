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
		Name:        "myproject",
		ModulePath:  "github.com/example/myproject",
		GoVersion:   "1.21",
		RepoURL:     "https://github.com/example/myproject",
		Description: "A great project",
		License:     "MIT",
	}

	got, err := renderer.Render("go_default.md", project)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	checks := []string{
		"# myproject",
		"A great project",
		"github.com/example/myproject",
		"Go 1.21",
		"https://github.com/example/myproject",
		"MIT",
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
