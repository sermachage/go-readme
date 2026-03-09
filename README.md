# go-readme

`go-readme` is a README automation CLI for Go modules. It parses project metadata
from `go.mod` and git, renders a template, and idempotently updates `README.md`
between managed markers so custom content is preserved.

If you are new to the project, start with `go-readme generate --dry-run` to see
what will be written before changing any files.

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

Generate from another directory:

```sh
go-readme generate --dir ./path/to/module
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
| `--dir` | `.` | Target project directory |
| `--description`, `-d` | empty | Project description |
| `--template`, `-t` | `go_default.md` | Embedded template name |
| `--dry-run` | `false` | Print output without writing README |
| `--force` | `false` | Overwrite entire README (skip marker replacement) |
| `--non-interactive` | `false` | Disable interactive prompt |

### Doctor Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--dir` | `.` | Target project directory to diagnose |

## What gets generated

- **Title / metadata** — module name, install command, go version
- **Repository** — git remote URL when configured
- **Description** — from flag or interactive prompt
- **License** — detected from license file name

## Managed markers (beginner-friendly)

`go-readme` updates only the auto-managed block in your README:

```md
<!-- go-readme:start -->
...generated content...
<!-- go-readme:end -->
```

Anything outside this block is your manual content and is not changed.

Backward compatibility:
- Older projects may still have legacy markers:
  `<!-- readmeaker:start -->` and `<!-- readmeaker:end -->`
- `go-readme` still reads those legacy markers and will safely migrate them to
  the new `go-readme` marker names on the next `generate` run.

## License

See [LICENSE](LICENSE).
