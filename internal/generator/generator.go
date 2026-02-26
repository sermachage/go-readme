// Package generator produces README.md content from analyzed Go package data.
package generator

import (
	"bytes"
	"fmt"
	"go/doc"
	"strings"
	"text/template"

	"github.com/sermachage/go-readme/internal/analyzer"
)

const readmeTmpl = `# {{.Name}}
{{- if .Doc}}

{{trimRight .Doc}}
{{- end}}

## Installation

` + "```" + `sh
go install {{.ImportPath}}@latest
` + "```" + `
{{- if .Examples}}

## Usage
{{range .Examples}}
` + "```" + `go
{{synopsis .Doc}}
` + "```" + `
{{- end}}
{{- end}}
{{- if .Funcs}}

## Functions
{{range .Funcs}}
### `+"`"+`{{.Name}}{{funcDecl .Decl}}`+"`"+`

{{trimRight .Doc}}
{{end}}
{{- end}}
{{- if .Types}}

## Types
{{range .Types}}
### `+"`"+`{{.Name}}`+"`"+`

{{trimRight .Doc}}
{{- if .Methods}}

**Methods**
{{range .Methods}}
- `+"`"+`{{.Name}}{{funcDecl .Decl}}`+"`"+` — {{oneLiner .Doc}}
{{- end}}
{{- end}}
{{end}}
{{- end}}
{{- if .Consts}}

## Constants
{{range .Consts}}
{{trimRight .Doc}}
` + "```" + `go
{{valueDecl .Decl}}
` + "```" + `
{{end}}
{{- end}}
{{- if .Vars}}

## Variables
{{range .Vars}}
{{trimRight .Doc}}
` + "```" + `go
{{valueDecl .Decl}}
` + "```" + `
{{end}}
{{- end}}
{{- if .HasLicense}}

## License

See [LICENSE](LICENSE).
{{- end}}
`

// Generate returns the README.md content for the given package.
func Generate(pkg *analyzer.Package) (string, error) {
	funcMap := template.FuncMap{
		"trimRight": func(s string) string { return strings.TrimRight(s, " \t\n\r") },
		"synopsis":  doc.Synopsis,
		"oneLiner":  oneLiner,
		"funcDecl":  funcDeclString,
		"valueDecl": valueDeclString,
	}

	tmpl, err := template.New("readme").Funcs(funcMap).Parse(readmeTmpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, pkg); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}
