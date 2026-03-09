# Changelog

All notable changes to this project are documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0](https://github.com/sermachage/go-readme/compare/v1.0.0...v1.1.0) (2026-03-09)


### Features

* add automated release workflow with Release Please and GoReleaser ([57be022](https://github.com/sermachage/go-readme/commit/57be022d956eaa2485c0c0f2ab60ebde0405d06c))
* migrate to go-readme markers with legacy compatibility ([0bbbfec](https://github.com/sermachage/go-readme/commit/0bbbfec9e2f2a4df33fda6f89e39ba07e60116a9))


### Bug Fixes

* align canonical go-readme CLI and harden checks/tests ([54721a3](https://github.com/sermachage/go-readme/commit/54721a35cd2a94a1d7f0a6912a4b098ac113123a))
* align canonical go-readme CLI and harden checks/tests ([466f211](https://github.com/sermachage/go-readme/commit/466f21166ced9457178ac9664192aa73e692f50f))
* restore require directives in go.mod after merge from main ([e2b1579](https://github.com/sermachage/go-readme/commit/e2b15791cdde7907b94b11ab6f2197a756ac547e))

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
