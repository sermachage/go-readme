package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sermachage/go-readme/internal/detectors"
	"github.com/sermachage/go-readme/internal/domain"
	"github.com/sermachage/go-readme/internal/parser"
)

func TestGenerate_ErrWhenNotGoModule(t *testing.T) {
	_, err := Generate(GenerateOptions{Dir: t.TempDir()})
	if err == nil {
		t.Fatal("expected error for directory without go.mod")
	}
}

func TestGenerate_DryRunDoesNotWriteReadme(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, "go.mod", "module github.com/example/dryrun\n\ngo 1.24\n")

	res, err := Generate(GenerateOptions{
		Dir:         dir,
		Description: "dry run project",
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if !res.Created {
		t.Fatal("Created = false, want true for missing README")
	}
	if !strings.Contains(res.Content, "<!-- readmeaker:start -->") {
		t.Fatal("generated content missing managed start marker")
	}
	if _, err := os.Stat(filepath.Join(dir, "README.md")); !os.IsNotExist(err) {
		t.Fatalf("README.md exists after dry run, err=%v", err)
	}
}

func TestGenerate_UpdatesMarkedSection(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, "go.mod", "module github.com/example/project\n\ngo 1.24\n")
	writeTestFile(t, dir, "README.md", strings.Join([]string{
		"# Project",
		"",
		"Intro text.",
		"",
		"<!-- readmeaker:start -->",
		"old generated content",
		"<!-- readmeaker:end -->",
		"",
		"Manual notes.",
	}, "\n"))

	res, err := Generate(GenerateOptions{
		Dir:         dir,
		Description: "new generated description",
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if res.Created {
		t.Fatal("Created = true, want false for existing README")
	}

	got := mustReadFile(t, filepath.Join(dir, "README.md"))
	if strings.Contains(got, "old generated content") {
		t.Fatal("expected old generated content to be replaced")
	}
	if !strings.Contains(got, "Manual notes.") {
		t.Fatal("manual README content after markers should be preserved")
	}
	if !strings.Contains(got, "new generated description") {
		t.Fatal("new generated content missing from README")
	}
}

func TestGenerate_ForceOverwritesEntireReadme(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, "go.mod", "module github.com/example/force\n\ngo 1.24\n")
	writeTestFile(t, dir, "README.md", "# Existing README\n\nManual section")

	res, err := Generate(GenerateOptions{
		Dir:         dir,
		Description: "force replace",
		Force:       true,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if res.Created {
		t.Fatal("Created = true, want false for existing README")
	}

	got := mustReadFile(t, filepath.Join(dir, "README.md"))
	if strings.Contains(got, "Manual section") {
		t.Fatal("manual content should be removed with Force=true")
	}
	if strings.Contains(got, "<!-- readmeaker:start -->") {
		t.Fatal("markers should not be used when Force=true")
	}
}

func TestGenerate_DetectsLicenseCandidates(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, "go.mod", "module github.com/example/license\n\ngo 1.24\n")
	writeTestFile(t, dir, "LICENSE.txt", "MIT")

	res, err := Generate(GenerateOptions{
		Dir:         dir,
		Description: "licensed project",
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if !strings.Contains(res.Content, "## License") {
		t.Fatal("expected rendered README to include license section")
	}
	if !strings.Contains(res.Content, "LICENSE.txt") {
		t.Fatal("expected rendered README to include detected license file name")
	}
}

func TestGenerateService_UsesInjectedDependencies(t *testing.T) {
	store := &fakeStore{}
	s := &GenerateService{
		Detector: fakeDetector{result: detectors.DetectionResult{IsGoProject: true}},
		GoMod: fakeGoModReader{
			info: &parser.GoModInfo{
				ModulePath: "github.com/example/injected",
				GoVersion:  "1.24",
			},
		},
		Git:      fakeGitReader{info: &parser.GitInfo{RemoteURL: "https://example.com/repo"}},
		Renderer: fakeRenderer{content: "generated content"},
		Store:    store,
	}

	res, err := s.Generate(GenerateOptions{
		Dir:    t.TempDir(),
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if store.writeCalls != 0 {
		t.Fatalf("expected no writes during dry-run, got %d", store.writeCalls)
	}
	if !strings.Contains(res.Content, "generated content") {
		t.Fatalf("unexpected result content: %q", res.Content)
	}
}

func TestGenerateService_PropagatesDependencyErrors(t *testing.T) {
	s := &GenerateService{
		Detector: fakeDetector{result: detectors.DetectionResult{IsGoProject: true}},
		GoMod:    fakeGoModReader{err: errors.New("gomod parse failed")},
		Git:      fakeGitReader{info: &parser.GitInfo{}},
		Renderer: fakeRenderer{content: "generated"},
		Store:    &fakeStore{},
	}

	_, err := s.Generate(GenerateOptions{Dir: t.TempDir()})
	if err == nil || !strings.Contains(err.Error(), "gomod parse failed") {
		t.Fatalf("expected gomod parse error, got: %v", err)
	}
}

func writeTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func mustReadFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

type fakeDetector struct {
	result detectors.DetectionResult
}

func (f fakeDetector) Detect(string) detectors.DetectionResult {
	return f.result
}

type fakeGoModReader struct {
	info *parser.GoModInfo
	err  error
}

func (f fakeGoModReader) ParseGoMod(string) (*parser.GoModInfo, error) {
	return f.info, f.err
}

type fakeGitReader struct {
	info *parser.GitInfo
}

func (f fakeGitReader) ParseGit(string) *parser.GitInfo {
	return f.info
}

type fakeRenderer struct {
	content string
	err     error
}

func (f fakeRenderer) Render(string, domain.Project) (string, error) {
	return f.content, f.err
}

type fakeStore struct {
	existing   string
	readErr    error
	writeErr   error
	writeCalls int
}

func (f *fakeStore) ReadExisting(string) (string, error) {
	return f.existing, f.readErr
}

func (f *fakeStore) Write(string, string) error {
	f.writeCalls++
	return f.writeErr
}
