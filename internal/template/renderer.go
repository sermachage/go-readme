// Package template provides README template rendering.
package template

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"github.com/sermachage/go-readme/internal/domain"
)

//go:embed templates/*.md
var templateFS embed.FS

// Renderer renders README templates against a Project.
type Renderer struct {
	fs embed.FS
}

// NewRenderer returns a Renderer backed by the embedded templates.
func NewRenderer() *Renderer {
	return &Renderer{fs: templateFS}
}

// Render renders the named template file with the given project data.
// templateName should be relative, e.g. "go_default.md".
func (r *Renderer) Render(templateName string, project domain.Project) (string, error) {
	path := "templates/" + templateName
	data, err := r.fs.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("loading template %q: %w", templateName, err)
	}

	tmpl, err := template.New(templateName).Parse(string(data))
	if err != nil {
		return "", fmt.Errorf("parsing template %q: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, project); err != nil {
		return "", fmt.Errorf("executing template %q: %w", templateName, err)
	}
	return buf.String(), nil
}
