package detectors_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/sermachage/go-readme/internal/detectors"
)

func testdataDir(subdir string) string {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..")
	return filepath.Join(root, "testdata", subdir)
}

func TestGoDetector_IsGoProject(t *testing.T) {
	d := &detectors.GoDetector{}
	result := d.Detect(testdataDir("valid_go_project"))
	if !result.IsGoProject {
		t.Error("expected IsGoProject=true for directory with go.mod")
	}
}

func TestGoDetector_NotGoProject(t *testing.T) {
	d := &detectors.GoDetector{}
	result := d.Detect(testdataDir("no_git"))
	if result.IsGoProject {
		t.Error("expected IsGoProject=false for directory without go.mod")
	}
}
