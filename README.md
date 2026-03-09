# go-readme

`go-readme` is a README automation CLI for Go modules. It parses project metadata
from `go.mod` and git, renders a template, and idempotently updates `README.md`
between managed markers so custom content is preserved.

## Installation

```sh
go install github.com/sermachage/go-readme@latest
```

## Usage

Run from inside a Go module to generate or update `README.md`:

```sh
go-readme generate
```

Preview without writing:

```sh
go-readme generate --dry-run
```

### Flags

| Command | Description |
|---------|-------------|
| `generate` | Generate or update README content |
| `doctor` | Check project setup (`go.mod`, git, remote, README) |
| `version` | Print CLI version |

### Generate Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--description`, `-d` | empty | Project description |
| `--template`, `-t` | `go_default.md` | Embedded template name |
| `--dry-run` | `false` | Print output without writing README |
| `--force` | `false` | Overwrite entire README (skip marker replacement) |
| `--non-interactive` | `false` | Disable interactive prompt |

## What gets generated

- **Title / metadata** — module name, install command, go version
- **Repository** — git remote URL when configured
- **Description** — from flag or interactive prompt
- **License** — detected from license file name

## License

See [LICENSE](LICENSE).
