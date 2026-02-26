package generator_test

import (
	"go/ast"
	"go/doc"
	"strings"
	"testing"

	"github.com/sermachage/go-readme/internal/analyzer"
	"github.com/sermachage/go-readme/internal/generator"
)

func TestGenerate_ContainsPackageName(t *testing.T) {
	pkg := &analyzer.Package{
		Name:       "mylib",
		ImportPath: "github.com/example/mylib",
		Doc:        "Package mylib does something useful.\n",
	}

	out, err := generator.Generate(pkg)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if !strings.Contains(out, "# mylib") {
		t.Errorf("output missing '# mylib'; got:\n%s", out)
	}
	if !strings.Contains(out, "go install github.com/example/mylib@latest") {
		t.Errorf("output missing installation instruction; got:\n%s", out)
	}
}

func TestGenerate_ContainsDoc(t *testing.T) {
	pkg := &analyzer.Package{
		Name:       "util",
		ImportPath: "github.com/example/util",
		Doc:        "Package util provides utility functions.\n",
	}

	out, err := generator.Generate(pkg)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if !strings.Contains(out, "Package util provides utility functions.") {
		t.Errorf("output missing package doc; got:\n%s", out)
	}
}

func TestGenerate_LicenseSection(t *testing.T) {
	pkg := &analyzer.Package{
		Name:       "lib",
		ImportPath: "example.com/lib",
		HasLicense: true,
	}

	out, err := generator.Generate(pkg)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if !strings.Contains(out, "## License") {
		t.Errorf("output missing License section; got:\n%s", out)
	}
	if !strings.Contains(out, "[LICENSE](LICENSE)") {
		t.Errorf("output missing LICENSE link; got:\n%s", out)
	}
}

func TestGenerate_NoLicenseSection(t *testing.T) {
	pkg := &analyzer.Package{
		Name:       "lib",
		ImportPath: "example.com/lib",
		HasLicense: false,
	}

	out, err := generator.Generate(pkg)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if strings.Contains(out, "## License") {
		t.Errorf("output contains License section but HasLicense=false; got:\n%s", out)
	}
}

func TestGenerate_TypesAndFunctions(t *testing.T) {
	// Build a minimal doc.Type and doc.Func using AST nodes.
	pkg := &analyzer.Package{
		Name:       "shapes",
		ImportPath: "example.com/shapes",
		Types: []*doc.Type{
			{
				Name: "Circle",
				Doc:  "Circle represents a circle.\n",
				Decl: &ast.GenDecl{},
			},
		},
		Funcs: []*doc.Func{
			{
				Name: "NewCircle",
				Doc:  "NewCircle creates a new Circle.\n",
				Decl: &ast.FuncDecl{
					Name: ast.NewIdent("NewCircle"),
					Type: &ast.FuncType{},
				},
			},
		},
	}

	out, err := generator.Generate(pkg)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if !strings.Contains(out, "## Types") {
		t.Errorf("output missing Types section; got:\n%s", out)
	}
	if !strings.Contains(out, "## Functions") {
		t.Errorf("output missing Functions section; got:\n%s", out)
	}
	if !strings.Contains(out, "Circle") {
		t.Errorf("output missing Circle type; got:\n%s", out)
	}
	if !strings.Contains(out, "NewCircle") {
		t.Errorf("output missing NewCircle function; got:\n%s", out)
	}
}
