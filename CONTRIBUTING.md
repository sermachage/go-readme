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

### Run tests

```sh
go test ./...
```

### Try the CLI locally

```sh
go install .
go-readme generate --dry-run
```

---

## Project layout

```
cmd/                    # Cobra commands (generate, doctor, version)
main.go                 # go-readme entrypoint
internal/
  app/                  # Application service orchestration
  parser/               # go.mod and git metadata parsing
  template/             # Embedded README templates
  markers/              # Idempotent managed block replacement
  writer/               # README read/write (atomic write path)
  gitmeta/              # Shared git metadata helpers
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
   go test ./...
   ```

5. **Open a pull request** against `main`. Fill in the PR template.

---

## Commit style

Use short, imperative present-tense commit messages. Examples:

```
add --dir flag to doctor command
fix marker migration for legacy readmeaker blocks
```

---

## Reporting bugs / requesting features

Please use the [issue templates](.github/ISSUE_TEMPLATE/) rather than opening a blank issue — it keeps triage fast.

---

## Code style

- Standard `gofmt` formatting is required (CI will catch it).
- Follow the conventions already established in the file you are editing.
- Export doc-comments are required on all exported identifiers.

## Documentation expectations (important)

- Keep README and CONTRIBUTING updates beginner-friendly.
- Prefer concrete examples over abstract descriptions.
- When behavior changes, include a short "what changed" note and a command example.

---

## License

By contributing you agree that your contributions will be licensed under the [MIT License](LICENSE).
