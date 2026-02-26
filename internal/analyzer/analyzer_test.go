package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sermachage/go-readme/internal/analyzer"
)

// writeFile is a helper that writes content to a file inside a temp dir.
func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("writeFile %s: %v", name, err)
	}
}

func TestAnalyze_BasicPackage(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "math.go", `// Package math provides basic math helpers.
package math

// Add returns the sum of a and b.
func Add(a, b int) int { return a + b }

// Sub returns the difference of a and b.
func Sub(a, b int) int { return a - b }
`)

	pkg, err := analyzer.Analyze(dir)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	if pkg.Name != "math" {
		t.Errorf("Name = %q, want %q", pkg.Name, "math")
	}
	if pkg.Doc == "" {
		t.Error("Doc is empty, want package-level comment")
	}
	if len(pkg.Funcs) != 2 {
		t.Errorf("Funcs count = %d, want 2", len(pkg.Funcs))
	}
}

func TestAnalyze_WithTypes(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "types.go", `// Package shapes defines geometric shapes.
package shapes

// Circle represents a circle.
type Circle struct {
	Radius float64
}

// Area returns the area of the circle.
func (c Circle) Area() float64 { return 3.14 * c.Radius * c.Radius }
`)

	pkg, err := analyzer.Analyze(dir)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	if len(pkg.Types) != 1 {
		t.Fatalf("Types count = %d, want 1", len(pkg.Types))
	}
	if pkg.Types[0].Name != "Circle" {
		t.Errorf("Types[0].Name = %q, want %q", pkg.Types[0].Name, "Circle")
	}
}

func TestAnalyze_WithConsts(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "consts.go", `// Package flags defines feature flags.
package flags

// MaxRetry is the maximum number of retry attempts.
const MaxRetry = 3
`)

	pkg, err := analyzer.Analyze(dir)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	if len(pkg.Consts) != 1 {
		t.Fatalf("Consts count = %d, want 1", len(pkg.Consts))
	}
}

func TestAnalyze_LicenseDetection(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "doc.go", `// Package example is an example.
package example
`)

	// Without LICENSE file.
	pkg, err := analyzer.Analyze(dir)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if pkg.HasLicense {
		t.Error("HasLicense = true, want false (no LICENSE file)")
	}

	// With LICENSE file.
	writeFile(t, dir, "LICENSE", "MIT License")
	pkg, err = analyzer.Analyze(dir)
	if err != nil {
		t.Fatalf("Analyze with LICENSE: %v", err)
	}
	if !pkg.HasLicense {
		t.Error("HasLicense = false, want true")
	}
}

func TestAnalyze_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	pkg, err := analyzer.Analyze(dir)
	if err != nil {
		t.Fatalf("Analyze empty dir: %v", err)
	}
	if pkg == nil {
		t.Fatal("Analyze returned nil Package")
	}
}
