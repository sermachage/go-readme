# go-readme

`go-readme` is the definitive README automation tool for Go. It inspects any Go package using the standard `go/doc` toolchain and produces a well-structured `README.md` — no configuration required.

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

Run from inside a Go module to generate a `README.md` for the current directory:

```sh
go-readme
```

Or point it at a specific package directory and output file:

```sh
go-readme -dir ./mypackage -output README.md
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-dir` | `.` | Directory of the Go package to document |
| `-output` | `README.md` | Output file path |

## What gets generated

- **Title** — the package name
- **Package doc** — the package-level documentation comment
- **Installation** — `go install` command with the correct import path (auto-detected from `go.mod`)
- **Functions** — all exported functions with their signatures and docs
- **Types** — all exported types with their docs and method list
- **Constants / Variables** — exported const and var blocks with docs
- **License** — link to `LICENSE` file if one exists in the package directory

## License

See [LICENSE](LICENSE).
