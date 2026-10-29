package disk

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUsageSumsApparentSizes(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "sub"), 0o755)
	os.WriteFile(filepath.Join(dir, "a.txt"), make([]byte, 4096), 0o644)
	os.WriteFile(filepath.Join(dir, "sub", "b.bin"), make([]byte, 1500), 0o644)
	var out bytes.Buffer
	if e := Usage(&out, []string{filepath.Join(dir, "a.txt"), filepath.Join(dir, "sub")}); e != nil {
		t.Fatal(e)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %q", out.String())
	}
	if !strings.HasSuffix(lines[0], filepath.Join(dir, "a.txt")) || !strings.Contains(lines[0], "4.0K") {
		t.Fatalf("wrong report for a.txt: %q", lines[0])
	}
	if !strings.Contains(lines[1], "1.5K") {
		t.Fatalf("wrong report for sub: %q", lines[1])
	}
}

func TestUsageDefaultListsHiddenThenVisible(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PWD", dir) // irrelevant; chdir matters
	os.WriteFile(filepath.Join(dir, ".hidden"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(dir, "visible"), []byte("x"), 0o644)
	old, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(old)
	var out bytes.Buffer
	if e := Usage(&out, nil); e != nil {
		t.Fatal(e)
	}
	s := out.String()
	hi := strings.Index(s, ".hidden")
	vi := strings.Index(s, "visible")
	if hi < 0 || vi < 0 || hi > vi {
		t.Fatalf("hidden files must list first: %q", s)
	}
}

func TestHumanize(t *testing.T) {
	cases := map[int64]string{512: "512", 4096: "4.0K", 1536: "1.5K", 12288: "12K", 1 << 30: "1.0G"}
	for in, want := range cases {
		if got := humanize(in); got != want {
			t.Fatalf("humanize(%d) = %s, want %s", in, got, want)
		}
	}
}
