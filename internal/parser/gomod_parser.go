// Package parser provides parsers for Go project metadata sources.
package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// GoModInfo holds information parsed from a go.mod file.
type GoModInfo struct {
	ModulePath   string
	GoVersion    string
	Dependencies []string
}

// ParseGoMod reads and parses the go.mod file in the given directory.
func ParseGoMod(dir string) (*GoModInfo, error) {
	path := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading go.mod: %w", err)
	}
	f, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, fmt.Errorf("parsing go.mod: %w", err)
	}

	info := &GoModInfo{
		ModulePath: f.Module.Mod.Path,
	}
	if f.Go != nil {
		info.GoVersion = f.Go.Version
	}
	for _, req := range f.Require {
		if !req.Indirect {
			info.Dependencies = append(info.Dependencies, req.Mod.Path)
		}
	}
	return info, nil
}
