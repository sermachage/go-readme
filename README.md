# go-readme

`go-readme` is a README automation CLI for Go modules. It parses project metadata
from `go.mod`, git, Go source, and common repo files, renders a template, and
idempotently updates `README.md` between managed markers so custom content is
preserved.

If you are new to the project, start with `go-readme generate --dry-run` to see
what will be written before changing any files.

## Installation

`go-readme` requires Go 1.24 or newer.

```sh
go install -v github.com/sermachage/go-readme/cmd/go-readme@latest
```

The `-v` flag prints the packages being compiled so you can see the install
progress. When the command returns to your shell prompt, the installation is
complete.

> **`go-readme` not found?** Make sure Go's binary directory is in your `PATH`:
>
> ```sh
> export PATH="$PATH:$(go env GOPATH)/bin"
> ```
>
> Add that line to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) to make it
> permanent.

> **Install failed with `permission denied` or `read-only file system`?**
> `go install` writes the binary to your Go bin directory. If the default
> location is not writable, install to a writable directory instead:
>
> ```sh
> mkdir -p "$HOME/.local/bin"
> GOBIN="$HOME/.local/bin" go install -v github.com/sermachage/go-readme/cmd/go-readme@latest
> export PATH="$HOME/.local/bin:$PATH"
> ```

Verify the installation:

```sh
go-readme version
```

## Usage

Run from inside a Go module to generate or update `README.md`:

```sh
go-readme generate
```

In interactive mode, `go-readme` now asks for a short project description, key
features, a usage example, configuration notes, and contributing notes when
those values were not already supplied by flags.

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
| `--features` | empty | Comma-separated key project features |
| `--usage-example` | empty | Usage command or code snippet |
| `--configuration` | empty | Configuration notes |
| `--contributing-notes` | empty | Contributor guidance |
| `--template`, `-t` | `go_default.md` | Embedded template name |
| `--dry-run` | `false` | Print output without writing README |
| `--force` | `false` | Overwrite entire README (skip marker replacement) |
| `--non-interactive` | `false` | Disable the interactive questionnaire |

### Doctor Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--dir` | `.` | Target project directory to diagnose |

## What gets generated

- **Overview** — project name plus description from a flag, prompt, or package
  doc comment fallback
- **Features** — optional bullet list from prompt or `--features`
- **Installation / usage** — install command plus a smarter default usage
  example that can be overridden
- **Configuration / development** — optional config notes plus standard Go
  build and test commands
- **Requirements / dependencies** — Go version and direct dependencies from
  `go.mod`
- **Repository / contributing / security** — git remote plus links to
  `CONTRIBUTING.md` and `SECURITY.md` when present
- **License** — detected from the license file name and linked in the output

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
