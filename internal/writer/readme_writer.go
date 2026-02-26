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
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing README: %w", err)
	}
	return nil
}
