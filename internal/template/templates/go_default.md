# {{ .Name }}

{{ .Description }}

## Installation

```bash
go install {{ .ModulePath }}@latest
```

## Usage

```go
import "{{ .ModulePath }}"
```

## Requirements

* Go {{ .GoVersion }}
{{- if .RepoURL }}

## Repository

[{{ .RepoURL }}]({{ .RepoURL }})
{{- end }}
{{- if .License }}

## License

{{ .License }}
{{- end }}
