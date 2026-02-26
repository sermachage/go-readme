// Package domain defines the core data models used throughout readmeaker.
package domain

// Project holds all metadata extracted from a Go project.
type Project struct {
	Name        string
	ModulePath  string
	GoVersion   string
	RepoURL     string
	Description string
	License     string
}
