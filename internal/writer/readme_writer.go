// Package writer handles reading and writing README files.
package writer

import (
	"fmt"
	"os"
	"path/filepath"
)

const readmeFile = "README.md"

// ReadExisting reads the existing README in dir, returning "" if it doesn't exist.
func ReadExisting(dir string) (string, error) {
	path := filepath.Join(dir, readmeFile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("reading README: %w", err)
	}
	return string(data), nil
}

// Write writes content to README.md in dir.
func Write(dir, content string) error {
	path := filepath.Join(dir, readmeFile)
	tmp, err := os.CreateTemp(dir, readmeFile+".tmp-*")
	if err != nil {
		return fmt.Errorf("writing README: create temp file: %w", err)
	}

	tmpPath := tmp.Name()
	cleanup := func() {
		_ = os.Remove(tmpPath)
	}

	if _, err := tmp.WriteString(content); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("writing README: write temp file: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("writing README: sync temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return fmt.Errorf("writing README: close temp file: %w", err)
	}
	if err := os.Chmod(tmpPath, 0o644); err != nil {
		cleanup()
		return fmt.Errorf("writing README: chmod temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		cleanup()
		return fmt.Errorf("writing README: rename temp file: %w", err)
	}

	if err := syncDir(dir); err != nil {
		return fmt.Errorf("writing README: %w", err)
	}
	return nil
}

func syncDir(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()
	return f.Sync()
}
