// Package analyzer parses Go source files and extracts documentation.
package analyzer

import (
	"bufio"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Package holds all documentation extracted from a Go package.
type Package struct {
	// Name is the package name.
	Name string
	// ImportPath is the module-relative import path, e.g. "github.com/org/repo".
	ImportPath string
	// Doc is the package-level documentation comment.
	Doc string
	// Funcs are the exported top-level functions.
	Funcs []*doc.Func
	// Types are the exported types, each carrying their own methods.
	Types []*doc.Type
	// Consts are the exported constant groups.
	Consts []*doc.Value
	// Vars are the exported variable groups.
	Vars []*doc.Value
	// Examples are the example functions.
	Examples []*doc.Example
	// HasLicense is true when a LICENSE file was found in the package directory.
	HasLicense bool
	// Dir is the absolute path of the analyzed directory.
	Dir string
}

// Analyze parses the Go source files in dir and returns a Package.
func Analyze(dir string) (*Package, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, absDir, func(fi os.FileInfo) bool {
		name := fi.Name()
		// Skip test files so they don't pollute the public API surface.
		return filepath.Ext(name) == ".go" && !strings.HasSuffix(name, "_test.go")
	}, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Find the first non-test package.
	importPath := detectImportPath(absDir)
	var dpkg *doc.Package
	for name, apkg := range pkgs {
		if len(name) > 5 && name[len(name)-5:] == "_test" {
			continue
		}
		files := mapToSlice(apkg.Files)
		dpkg, err = doc.NewFromFiles(fset, files, importPath)
		if err != nil {
			return nil, err
		}
		break
	}

	if dpkg == nil {
		// Fall back to the first package found (handles _test-only scenarios).
		for _, apkg := range pkgs {
			files := mapToSlice(apkg.Files)
			dpkg, err = doc.NewFromFiles(fset, files, importPath)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	if dpkg == nil {
		return &Package{Dir: absDir}, nil
	}

	_, licErr := os.Stat(filepath.Join(absDir, "LICENSE"))
	hasLicense := licErr == nil

	return &Package{
		Name:       dpkg.Name,
		ImportPath: dpkg.ImportPath,
		Doc:        dpkg.Doc,
		Funcs:      dpkg.Funcs,
		Types:      dpkg.Types,
		Consts:     dpkg.Consts,
		Vars:       dpkg.Vars,
		Examples:   dpkg.Examples,
		HasLicense: hasLicense,
		Dir:        absDir,
	}, nil
}

// mapToSlice converts map[string]*ast.File to []*ast.File.
func mapToSlice(m map[string]*ast.File) []*ast.File {
	files := make([]*ast.File, 0, len(m))
	for _, f := range m {
		files = append(files, f)
	}
	return files
}

// detectImportPath walks up from dir looking for a go.mod file. When found it
// computes the import path of dir relative to the module root.
func detectImportPath(dir string) string {
	modPath, modDir := findGoMod(dir)
	if modPath == "" {
		return ""
	}
	rel, err := filepath.Rel(modDir, dir)
	if err != nil || rel == "." {
		return modPath
	}
	return modPath + "/" + filepath.ToSlash(rel)
}

// findGoMod searches for go.mod starting at dir and moving to parent
// directories. It returns the module path declared in go.mod and the directory
// that contains go.mod.
func findGoMod(dir string) (modulePath, moduleDir string) {
	current := dir
	for {
		candidate := filepath.Join(current, "go.mod")
		f, err := os.Open(candidate)
		if err == nil {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if strings.HasPrefix(line, "module ") {
					_ = f.Close()
					return strings.TrimSpace(strings.TrimPrefix(line, "module")), current
				}
			}
			_ = f.Close()
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", ""
}
