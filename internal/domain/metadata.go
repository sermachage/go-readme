// Package domain defines the core data models used throughout readmeaker.
package domain

// Metadata holds raw extracted values before they are assembled into a Project.
type Metadata struct {
	ModulePath   string
	GoVersion    string
	Dependencies []string
	RepoURL      string
	Branch       string
}
