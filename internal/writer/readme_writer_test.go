package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadExisting_MissingReadmeReturnsEmpty(t *testing.T) {
	got, err := ReadExisting(t.TempDir())
	if err != nil {
		t.Fatalf("ReadExisting: %v", err)
	}
	if got != "" {
		t.Fatalf("ReadExisting = %q, want empty string", got)
	}
}

func TestWriteAndReadExisting(t *testing.T) {
	dir := t.TempDir()
	want := "# Title\n\ncontent"
	if err := Write(dir, want); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := ReadExisting(dir)
	if err != nil {
		t.Fatalf("ReadExisting: %v", err)
	}
	if got != want {
		t.Fatalf("ReadExisting = %q, want %q", got, want)
	}
}

func TestReadExisting_WrapsReadError(t *testing.T) {
	dir := t.TempDir()
	notDir := filepath.Join(dir, "plain-file")
	if err := os.WriteFile(notDir, []byte("x"), 0o644); err != nil {
		t.Fatalf("write setup file: %v", err)
	}

	_, err := ReadExisting(notDir)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "reading README") {
		t.Fatalf("expected wrapped error, got: %v", err)
	}
}

func TestWrite_WrapsWriteError(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "does-not-exist")
	err := Write(missing, "content")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "writing README") {
		t.Fatalf("expected wrapped error, got: %v", err)
	}
}
