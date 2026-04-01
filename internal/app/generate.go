// Package app implements the application service layer.
package app

import (
	"fmt"
	"go/doc"
	"os"
	"path/filepath"
	"strings"

	"github.com/sermachage/go-readme/internal/analyzer"
	"github.com/sermachage/go-readme/internal/detectors"
	"github.com/sermachage/go-readme/internal/domain"
	"github.com/sermachage/go-readme/internal/markers"
	"github.com/sermachage/go-readme/internal/parser"
	tmpl "github.com/sermachage/go-readme/internal/template"
	"github.com/sermachage/go-readme/internal/writer"
)

// ProjectDetector determines whether a target directory is a supported project.
type ProjectDetector interface {
	Detect(dir string) detectors.DetectionResult
}

// GoModReader parses go.mod metadata.
type GoModReader interface {
	ParseGoMod(dir string) (*parser.GoModInfo, error)
}

// GitReader parses git metadata.
type GitReader interface {
	ParseGit(dir string) *parser.GitInfo
}

// SourceAnalyzer extracts lightweight documentation metadata from Go source.
type SourceAnalyzer interface {
	Analyze(dir string) (*analyzer.Package, error)
}

// ProjectRenderer renders a README template from project metadata.
type ProjectRenderer interface {
	Render(templateName string, project domain.Project) (string, error)
}

// ReadmeStore abstracts reading and writing README files.
type ReadmeStore interface {
	ReadExisting(dir string) (string, error)
	Write(dir, content string) error
}

// GenerateService orchestrates README generation dependencies.
type GenerateService struct {
	Detector ProjectDetector
	GoMod    GoModReader
	Git      GitReader
	Source   SourceAnalyzer
	Renderer ProjectRenderer
	Store    ReadmeStore
}

// GenerateOptions controls the behaviour of the Generate service.
type GenerateOptions struct {
	// Dir is the target project directory. Defaults to the current directory.
	Dir string
	// Description is an optional description to embed in the README.
	Description string
	// Features is an optional list of key project features.
	Features []string
	// UsageExample is an optional usage snippet or command.
	UsageExample string
	// Configuration is optional configuration guidance.
	Configuration string
	// Contributing is optional contributor guidance.
	Contributing string
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

type defaultGoModReader struct{}

func (r defaultGoModReader) ParseGoMod(dir string) (*parser.GoModInfo, error) {
	return parser.ParseGoMod(dir)
}

type defaultGitReader struct{}

func (r defaultGitReader) ParseGit(dir string) *parser.GitInfo {
	return parser.ParseGit(dir)
}

type defaultSourceAnalyzer struct{}

func (a defaultSourceAnalyzer) Analyze(dir string) (*analyzer.Package, error) {
	return analyzer.Analyze(dir)
}

type defaultReadmeStore struct{}

func (s defaultReadmeStore) ReadExisting(dir string) (string, error) {
	return writer.ReadExisting(dir)
}

func (s defaultReadmeStore) Write(dir, content string) error {
	return writer.Write(dir, content)
}

// NewGenerateService builds the default README generation service.
func NewGenerateService() *GenerateService {
	return &GenerateService{
		Detector: &detectors.GoDetector{},
		GoMod:    defaultGoModReader{},
		Git:      defaultGitReader{},
		Source:   defaultSourceAnalyzer{},
		Renderer: tmpl.NewRenderer(),
		Store:    defaultReadmeStore{},
	}
}

// Generate detects the project, extracts metadata, renders a template, and
// writes (or updates) README.md.
func Generate(opts GenerateOptions) (*GenerateResult, error) {
	return NewGenerateService().Generate(opts)
}

// Generate detects the project, extracts metadata, renders a template, and
// writes (or updates) README.md.
func (s *GenerateService) Generate(opts GenerateOptions) (*GenerateResult, error) {
	if opts.Dir == "" {
		opts.Dir = "."
	}
	if opts.Template == "" {
		opts.Template = "go_default.md"
	}

	result := s.Detector.Detect(opts.Dir)
	if !result.IsGoProject {
		return nil, fmt.Errorf("no go.mod found in %s – not a Go module project", opts.Dir)
	}

	gomod, err := s.GoMod.ParseGoMod(opts.Dir)
	if err != nil {
		return nil, err
	}

	git := s.Git.ParseGit(opts.Dir)
	var pkg *analyzer.Package
	if s.Source != nil {
		analyzed, err := s.Source.Analyze(opts.Dir)
		if err == nil {
			pkg = analyzed
		}
	}

	features := cleanList(opts.Features)
	dependencies, additionalDependencies := summarizeDependencies(gomod.Dependencies, 8)
	usageExample, usageLanguage := resolveUsageExample(opts.UsageExample, gomod.ModulePath, pkg)

	project := domain.Project{
		Name:                   moduleName(gomod.ModulePath),
		ModulePath:             gomod.ModulePath,
		GoVersion:              gomod.GoVersion,
		RepoURL:                git.RemoteURL,
		Description:            resolveDescription(opts.Description, pkg),
		Features:               features,
		UsageExample:           usageExample,
		UsageLanguage:          usageLanguage,
		Configuration:          strings.TrimSpace(opts.Configuration),
		Dependencies:           dependencies,
		AdditionalDependencies: additionalDependencies,
		Contributing:           strings.TrimSpace(opts.Contributing),
		ContributingGuide:      detectContributingGuide(opts.Dir),
		SecurityPolicy:         detectSecurityPolicy(opts.Dir),
		License:                detectLicense(opts.Dir),
	}

	rendered, err := s.Renderer.Render(opts.Template, project)
	if err != nil {
		return nil, err
	}

	existing, err := s.Store.ReadExisting(opts.Dir)
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
		if err := s.Store.Write(opts.Dir, content); err != nil {
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

func resolveDescription(description string, pkg *analyzer.Package) string {
	if trimmed := strings.TrimSpace(description); trimmed != "" {
		return trimmed
	}
	if pkg == nil {
		return ""
	}
	return strings.TrimSpace(doc.Synopsis(pkg.Doc))
}

func resolveUsageExample(usageExample, modulePath string, pkg *analyzer.Package) (string, string) {
	if trimmed := strings.TrimSpace(usageExample); trimmed != "" {
		return trimmed, inferCodeFenceLanguage(trimmed)
	}
	return defaultUsageExample(modulePath, pkg)
}

func defaultUsageExample(modulePath string, pkg *analyzer.Package) (string, string) {
	if pkg != nil && pkg.Name == "main" {
		return fmt.Sprintf("%s --help", moduleName(modulePath)), "sh"
	}
	return fmt.Sprintf("import %q", modulePath), "go"
}

func inferCodeFenceLanguage(example string) string {
	trimmed := strings.TrimSpace(example)
	switch {
	case trimmed == "":
		return "text"
	case strings.Contains(trimmed, "package ") || strings.Contains(trimmed, "import ") || strings.Contains(trimmed, "func "):
		return "go"
	case strings.HasPrefix(trimmed, "$ "):
		return "sh"
	case strings.HasPrefix(trimmed, "go ") || strings.HasPrefix(trimmed, "./") || strings.HasPrefix(trimmed, "make "):
		return "sh"
	default:
		return "text"
	}
}

func cleanList(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(values))
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		parts := strings.FieldsFunc(value, func(r rune) bool {
			return r == ',' || r == '\n'
		})
		for _, part := range parts {
			item := strings.TrimSpace(part)
			if item == "" {
				continue
			}
			if _, ok := seen[item]; ok {
				continue
			}
			seen[item] = struct{}{}
			cleaned = append(cleaned, item)
		}
	}
	if len(cleaned) == 0 {
		return nil
	}
	return cleaned
}

func summarizeDependencies(deps []string, limit int) ([]string, int) {
	cleaned := cleanList(deps)
	if len(cleaned) <= limit {
		return cleaned, 0
	}
	return cleaned[:limit], len(cleaned) - limit
}

// detectLicense looks for a LICENSE file and returns its name, or "".
func detectLicense(dir string) string {
	return detectDocFile(dir, []string{"LICENSE", "LICENSE.md", "LICENSE.txt", "LICENCE"})
}

func detectContributingGuide(dir string) string {
	return detectDocFile(dir, []string{"CONTRIBUTING.md", "CONTRIBUTING.txt", "CONTRIBUTING"})
}

func detectSecurityPolicy(dir string) string {
	return detectDocFile(dir, []string{"SECURITY.md", "SECURITY.txt", "SECURITY"})
}

func detectDocFile(dir string, candidates []string) string {
	for _, name := range candidates {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return name
		}
	}
	return ""
}
