# Changelog

All notable changes to this project are documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.0.0 (2026-02-26)


### Features

* add automated release workflow with Release Please and GoReleaser ([57be022](https://github.com/sermachage/go-readme/commit/57be022d956eaa2485c0c0f2ab60ebde0405d06c))


### Bug Fixes

* restore require directives in go.mod after merge from main ([e2b1579](https://github.com/sermachage/go-readme/commit/e2b15791cdde7907b94b11ab6f2197a756ac547e))

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
