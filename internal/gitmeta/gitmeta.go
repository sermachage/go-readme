// Package gitmeta provides shared git metadata utilities.
package gitmeta

import (
	"os/exec"
	"strings"
)

// Info holds git metadata extracted from a repository.
type Info struct {
	RemoteURL string
	Branch    string
}

// Parse extracts git metadata from the repository rooted at dir.
// Non-fatal errors return empty fields.
func Parse(dir string) *Info {
	info := &Info{}

	remote, err := output(dir, "config", "--get", "remote.origin.url")
	if err == nil {
		info.RemoteURL = NormalizeURL(remote)
	}

	branch, err := output(dir, "branch", "--show-current")
	if err == nil {
		info.Branch = strings.TrimSpace(branch)
	}

	return info
}

// IsRepository reports whether dir is inside a git work tree.
func IsRepository(dir string) (bool, error) {
	out, err := output(dir, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) == "true", nil
}

// RemoteOriginURL returns the normalized origin URL.
func RemoteOriginURL(dir string) (string, error) {
	out, err := output(dir, "config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}
	return NormalizeURL(out), nil
}

// NormalizeURL converts SSH remote URLs to HTTPS form and trims .git suffix.
func NormalizeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "git@") {
		raw = strings.TrimPrefix(raw, "git@")
		raw = strings.Replace(raw, ":", "/", 1)
		raw = "https://" + raw
	}
	return strings.TrimSuffix(raw, ".git")
}

func output(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
