package parser_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/sermachage/go-readme/internal/parser"
)

func testdataDir(t *testing.T, subdir string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file location")
	}
	// file is in internal/parser/, testdata is two levels up
	root := filepath.Join(filepath.Dir(file), "..", "..")
	return filepath.Join(root, "testdata", subdir)
}

func TestParseGoMod_Valid(t *testing.T) {
	dir := testdataDir(t, "valid_go_project")
	info, err := parser.ParseGoMod(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ModulePath != "github.com/example/myproject" {
		t.Errorf("ModulePath = %q, want %q", info.ModulePath, "github.com/example/myproject")
	}
	if info.GoVersion != "1.21" {
		t.Errorf("GoVersion = %q, want %q", info.GoVersion, "1.21")
	}
	// Only direct dependencies should be listed
	wantDeps := []string{"github.com/spf13/cobra", "github.com/some/dep"}
	if len(info.Dependencies) != len(wantDeps) {
		t.Errorf("len(Dependencies) = %d, want %d; got %v", len(info.Dependencies), len(wantDeps), info.Dependencies)
	}
}

func TestParseGoMod_NoFile(t *testing.T) {
	_, err := parser.ParseGoMod(testdataDir(t, "no_git"))
	if err == nil {
		t.Fatal("expected error for missing go.mod, got nil")
	}
}

func TestParseGoMod_Malformed(t *testing.T) {
	_, err := parser.ParseGoMod(testdataDir(t, "malformed_gomod"))
	if err == nil {
		t.Fatal("expected error for malformed go.mod, got nil")
	}
}
