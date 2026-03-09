package gitmeta

import (
	"os/exec"
	"strings"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "https://github.com/user/repo.git", want: "https://github.com/user/repo"},
		{in: "https://github.com/user/repo", want: "https://github.com/user/repo"},
		{in: "git@github.com:user/repo.git", want: "https://github.com/user/repo"},
		{in: "  git@github.com:user/repo  ", want: "https://github.com/user/repo"},
		{in: "", want: ""},
	}

	for _, tc := range tests {
		if got := NormalizeURL(tc.in); got != tc.want {
			t.Fatalf("NormalizeURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestIsRepository(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")

	ok, err := IsRepository(dir)
	if err != nil {
		t.Fatalf("IsRepository: %v", err)
	}
	if !ok {
		t.Fatal("expected IsRepository to return true")
	}
}

func TestRemoteOriginURL(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")
	runGit(t, dir, "remote", "add", "origin", "git@github.com:example/project.git")

	got, err := RemoteOriginURL(dir)
	if err != nil {
		t.Fatalf("RemoteOriginURL: %v", err)
	}
	if got != "https://github.com/example/project" {
		t.Fatalf("unexpected normalized remote: %q", got)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(out))
	}
}
