// Package detectors provides project-type detection logic.
package detectors

import (
	"os"
	"path/filepath"
)

// DetectionResult holds the outcome of a project detection run.
type DetectionResult struct {
	IsGoProject bool
	ModuleName  string
}

// GoDetector detects whether a directory is a Go module project.
type GoDetector struct{}

// Detect checks the given directory for a go.mod file and returns a
// DetectionResult.
func (d *GoDetector) Detect(dir string) DetectionResult {
	gomod := filepath.Join(dir, "go.mod")
	_, err := os.Stat(gomod)
	if err != nil {
		return DetectionResult{}
	}
	return DetectionResult{IsGoProject: true}
}
