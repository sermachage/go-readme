# Contributing to go-readme

Thank you for your interest in contributing! The guidelines below will help you get set up quickly and ensure your pull request is reviewed smoothly.

---

## Getting started

### Prerequisites

| Tool | Minimum version |
|------|----------------|
| [Go](https://go.dev/dl/) | 1.22 |
| Git | any recent version |

### Clone and build

```sh
git clone https://github.com/sermachage/go-readme.git
cd go-readme
go build ./...
```

### Run the tests

```sh
go test -race ./...
```

### Install locally for manual testing

```sh
go install ./cmd/go-readme
go-readme -dir ./internal/analyzer
```

---

## Project layout

```
cmd/go-readme/          # CLI entry point
internal/
  analyzer/             # Go source → Package struct (uses go/doc)
  generator/            # Package struct → README.md (uses text/template)
```

---

## Making changes

1. **Fork** the repository and create a branch from `main`:

   ```sh
   git checkout -b feat/my-feature
   ```

2. **Write code.** Keep changes focused; one logical change per PR.

3. **Add or update tests** in the relevant `*_test.go` file. All tests live alongside the code they test.

4. **Verify** before pushing:

   ```sh
   go build ./...
   go vet ./...
   go test -race ./...
   ```

5. **Open a pull request** against `main`. Fill in the PR template.

---

## Commit style

Use short, imperative present-tense commit messages:

```
add -stdout flag to print README to console
fix import-path detection for nested modules
```

---

## Reporting bugs / requesting features

Please use the [issue templates](.github/ISSUE_TEMPLATE/) rather than opening a blank issue — it keeps triage fast.

---

## Code style

- Standard `gofmt` formatting is required (CI will catch it).
- Follow the conventions already established in the file you are editing.
- Export doc-comments are required on all exported identifiers.

---

## License

By contributing you agree that your contributions will be licensed under the [MIT License](LICENSE).
