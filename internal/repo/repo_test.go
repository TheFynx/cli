package repo

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckReportsReposAndSummary(t *testing.T) {
	if _, e := exec.LookPath("git"); e != nil {
		t.Skip("git unavailable")
	}
	base := t.TempDir()
	clean := filepath.Join(base, "clean")
	dirty := filepath.Join(base, "dirty")
	for _, d := range []string{clean, dirty} {
		os.MkdirAll(d, 0o755)
		git := func(args ...string) {
			c := exec.Command("git", append([]string{"-C", d, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)...)
			c.Env = append(os.Environ(), "GIT_AUTHOR_DATE=2020-01-01T00:00:00Z")
			if e := c.Run(); e != nil {
				t.Fatalf("git %v: %v", args, e)
			}
		}
		git("init", "-q")
		os.WriteFile(filepath.Join(d, "f.txt"), []byte("one"), 0o644)
		git("add", ".")
		git("commit", "-qm", "one")
	}
	// Dirty repo: uncommitted change.
	os.WriteFile(filepath.Join(dirty, "f.txt"), []byte("two"), 0o644)

	var out bytes.Buffer
	if e := Check(base, &out); e != nil {
		t.Fatal(e)
	}
	s := out.String()
	if !strings.Contains(s, "clean") || !strings.Contains(s, "dirty") {
		t.Fatalf("missing repos in report: %q", s)
	}
	if !strings.Contains(s, "2 repositories: 0 with unmerged branches, 1 with local changes") {
		t.Fatalf("wrong summary: %q", s)
	}
}

func TestContentsGroupsAndHeaders(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "cmd", "inner"), 0o755)
	os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n"), 0o644)
	os.WriteFile(filepath.Join(root, "readme.md"), []byte("skip\n"), 0o644)
	os.WriteFile(filepath.Join(root, "cmd", "a.go"), []byte("package cmd\n"), 0o644)
	os.WriteFile(filepath.Join(root, "cmd", "inner", "b.go"), []byte("package inner\n"), 0o644)
	os.WriteFile(filepath.Join(root, ".goreleaser.yaml"), []byte("builds: []\n"), 0o644)

	var out bytes.Buffer
	target := filepath.Join(root, "..", filepath.Base(root)+"-dump.txt")
	if e := Contents(root, target, &out); e != nil {
		t.Fatal(e)
	}
	b, e := os.ReadFile(target)
	if e != nil {
		t.Fatal(e)
	}
	s := string(b)
	for _, want := range []string{"package main\n", "\n---\n\n", ".goreleaser.yaml content marker", "# cmd", "# cmd/a.go", "# cmd/inner/b.go", "package cmd", "package inner"} {
		if want == ".goreleaser.yaml content marker" {
			if !strings.Contains(s, "builds: []") {
				t.Fatal("goreleaser section missing")
			}
			continue
		}
		if !strings.Contains(s, want) {
			t.Fatalf("dump missing %q:\n%s", want, s)
		}
	}
	if strings.Contains(s, "readme.md") {
		t.Fatal("non-code files must be skipped")
	}
	if !strings.Contains(out.String(), "saved to") {
		t.Fatalf("missing confirmation: %q", out.String())
	}
}
