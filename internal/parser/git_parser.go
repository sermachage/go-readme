// Package parser provides parsers for Go project metadata sources.
package parser

import (
	"github.com/sermachage/go-readme/internal/gitmeta"
)

// GitInfo holds information extracted from a git repository.
type GitInfo struct {
	RemoteURL string
	Branch    string
}

// ParseGit extracts git metadata from the repository rooted at dir.
// Non-fatal errors (e.g. no remote configured) result in empty fields.
func ParseGit(dir string) *GitInfo {
	info := gitmeta.Parse(dir)
	return &GitInfo{
		RemoteURL: info.RemoteURL,
		Branch:    info.Branch,
	}
}

// NormalizeGitURL converts SSH remote URLs to HTTPS form.
// e.g. git@github.com:user/repo.git → https://github.com/user/repo
func NormalizeGitURL(raw string) string {
	return gitmeta.NormalizeURL(raw)
}
