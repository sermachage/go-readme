// Package app implements the application service layer.
package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sermachage/go-readme/internal/detectors"
	"github.com/sermachage/go-readme/internal/domain"
	"github.com/sermachage/go-readme/internal/markers"
	"github.com/sermachage/go-readme/internal/parser"
	tmpl "github.com/sermachage/go-readme/internal/template"
	"github.com/sermachage/go-readme/internal/writer"
)

// GenerateOptions controls the behaviour of the Generate service.
type GenerateOptions struct {
	// Dir is the target project directory. Defaults to the current directory.
	Dir string
	// Description is an optional description to embed in the README.
	Description string
	// Template is the template file name (without path). Defaults to "go_default.md".
	Template string
	// DryRun prints the output without writing to disk.
	DryRun bool
	// Force overwrites an existing README entirely (ignores markers).
	Force bool
}

// GenerateResult contains the output of a Generate call.
type GenerateResult struct {
	Content    string
	OutputPath string
	Created    bool
}

// Generate detects the project, extracts metadata, renders a template, and
// writes (or updates) README.md.
func Generate(opts GenerateOptions) (*GenerateResult, error) {
	if opts.Dir == "" {
		opts.Dir = "."
	}
	if opts.Template == "" {
		opts.Template = "go_default.md"
	}

	detector := &detectors.GoDetector{}
	result := detector.Detect(opts.Dir)
	if !result.IsGoProject {
		return nil, fmt.Errorf("no go.mod found in %s – not a Go module project", opts.Dir)
	}

	gomod, err := parser.ParseGoMod(opts.Dir)
	if err != nil {
		return nil, err
	}

	git := parser.ParseGit(opts.Dir)

	project := domain.Project{
		Name:        moduleName(gomod.ModulePath),
		ModulePath:  gomod.ModulePath,
		GoVersion:   gomod.GoVersion,
		RepoURL:     git.RemoteURL,
		Description: opts.Description,
		License:     detectLicense(opts.Dir),
	}

	renderer := tmpl.NewRenderer()
	rendered, err := renderer.Render(opts.Template, project)
	if err != nil {
		return nil, err
	}

	existing, err := writer.ReadExisting(opts.Dir)
	if err != nil {
		return nil, err
	}

	var content string
	created := existing == ""

	if opts.Force {
		content = rendered
	} else {
		content = markers.Replace(existing, rendered)
	}

	res := &GenerateResult{
		Content:    content,
		OutputPath: filepath.Join(opts.Dir, "README.md"),
		Created:    created,
	}

	if !opts.DryRun {
		if err := writer.Write(opts.Dir, content); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// moduleName extracts the short project name from a module path.
// e.g. "github.com/user/myproject" → "myproject"
func moduleName(modulePath string) string {
	parts := strings.Split(modulePath, "/")
	return parts[len(parts)-1]
}

// detectLicense looks for a LICENSE file and returns its name, or "".
func detectLicense(dir string) string {
	candidates := []string{"LICENSE", "LICENSE.md", "LICENSE.txt", "LICENCE"}
	for _, name := range candidates {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return name
		}
	}
	return ""
}
