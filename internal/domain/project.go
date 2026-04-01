// Package domain defines the core data models used throughout go-readme.
package domain

// Project holds all metadata extracted from a Go project.
type Project struct {
	Name                   string
	ModulePath             string
	GoVersion              string
	RepoURL                string
	Description            string
	Features               []string
	UsageExample           string
	UsageLanguage          string
	Configuration          string
	Dependencies           []string
	AdditionalDependencies int
	Contributing           string
	ContributingGuide      string
	SecurityPolicy         string
	License                string
}
