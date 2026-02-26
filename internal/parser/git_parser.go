// Package parser provides parsers for Go project metadata sources.
package parser

import (
	"os/exec"
	"strings"
)

// GitInfo holds information extracted from a git repository.
type GitInfo struct {
	RemoteURL string
	Branch    string
}

// ParseGit extracts git metadata from the repository rooted at dir.
// Non-fatal errors (e.g. no remote configured) result in empty fields.
func ParseGit(dir string) *GitInfo {
	info := &GitInfo{}
	info.RemoteURL = normalizeURL(gitOutput(dir, "git", "config", "--get", "remote.origin.url"))
	info.Branch = gitOutput(dir, "git", "branch", "--show-current")
	return info
}

// NormalizeGitURL converts SSH remote URLs to HTTPS form.
// e.g. git@github.com:user/repo.git → https://github.com/user/repo
func NormalizeGitURL(raw string) string {
	return normalizeURL(raw)
}

func normalizeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	// SSH form: git@github.com:user/repo.git
	if strings.HasPrefix(raw, "git@") {
		raw = strings.TrimPrefix(raw, "git@")
		raw = strings.Replace(raw, ":", "/", 1)
		raw = "https://" + raw
	}
	raw = strings.TrimSuffix(raw, ".git")
	return raw
}

// gitOutput runs a git command in dir and returns trimmed stdout, or "" on error.
func gitOutput(dir string, name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
