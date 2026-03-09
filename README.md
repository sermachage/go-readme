# go-readme

`go-readme` is a README automation CLI for Go modules. It parses project metadata
from `go.mod` and git, renders a template, and idempotently updates `README.md`
between managed markers so custom content is preserved.

## Installation

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

Verify the installation:

```sh
go-readme version
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
