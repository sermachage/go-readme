# Changelog

All notable changes to this project are documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial implementation of `go-readme` CLI
- `internal/analyzer`: Go source parser using `go/doc`; auto-detects module import path from `go.mod`
- `internal/generator`: README template renderer using `text/template`
- `-dir` flag to point at any Go package directory (default `.`)
- `-output` flag to set the output file path (default `README.md`)
- Generated sections: title, package doc, installation, functions, types, constants, variables, license
- CI workflow for Go 1.22 / 1.23 / 1.24
- Release workflow via GoReleaser
- Issue templates and contribution guidelines
