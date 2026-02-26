package parser_test

import (
	"testing"

	"github.com/sermachage/go-readme/internal/parser"
)

func TestNormalizeGitURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "https with .git suffix",
			input: "https://github.com/user/repo.git",
			want:  "https://github.com/user/repo",
		},
		{
			name:  "https without .git suffix",
			input: "https://github.com/user/repo",
			want:  "https://github.com/user/repo",
		},
		{
			name:  "ssh form",
			input: "git@github.com:user/repo.git",
			want:  "https://github.com/user/repo",
		},
		{
			name:  "ssh form without .git",
			input: "git@github.com:user/repo",
			want:  "https://github.com/user/repo",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "whitespace",
			input: "  https://github.com/user/repo.git  ",
			want:  "https://github.com/user/repo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := parser.NormalizeGitURL(tc.input)
			if got != tc.want {
				t.Errorf("NormalizeGitURL(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
