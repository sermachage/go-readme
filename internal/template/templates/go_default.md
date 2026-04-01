# {{ .Name }}
{{ if .Description }}
{{ .Description }}
{{ end }}
{{ if .Features }}
## Features

{{- range .Features }}
- {{ . }}
{{- end }}

{{ end }}
## Installation

```bash
go install {{ .ModulePath }}@latest
```

## Usage

```{{ .UsageLanguage }}
{{ .UsageExample }}
```
{{ if .Configuration }}
## Configuration

{{ .Configuration }}
{{ end }}
## Development

```bash
go test ./...
go build ./...
```
{{ if .GoVersion }}
## Requirements

- Go {{ .GoVersion }}
{{ end }}
{{ if .Dependencies }}
## Dependencies

{{- range .Dependencies }}
- `{{ . }}`
{{- end }}
{{ if gt .AdditionalDependencies 0 }}
- ...and {{ .AdditionalDependencies }} more direct dependencies in `go.mod`
{{ end }}

{{ end }}
{{ if .RepoURL }}
## Repository

[{{ .RepoURL }}]({{ .RepoURL }})
{{ end }}
{{ if or .Contributing .ContributingGuide }}
## Contributing

{{ if .Contributing }}
{{ .Contributing }}
{{ if .ContributingGuide }}

See [{{ .ContributingGuide }}]({{ .ContributingGuide }}) for the full contribution guide.
{{ end }}
{{ else }}
Contributions are welcome. See [{{ .ContributingGuide }}]({{ .ContributingGuide }}) for guidelines.
{{ end }}
{{ end }}
{{ if .SecurityPolicy }}
## Security

Review [{{ .SecurityPolicy }}]({{ .SecurityPolicy }}) before reporting vulnerabilities.
{{ end }}
{{ if .License }}
## License

See [{{ .License }}]({{ .License }}).
{{ end }}
