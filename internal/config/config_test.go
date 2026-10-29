package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingFileIsDefaults(t *testing.T) {
	t.Setenv("FYNX_CONFIG", filepath.Join(t.TempDir(), "absent.toml"))
	c, e := Load()
	if e != nil {
		t.Fatal(e)
	}
	if c.Upscale.Width != 0 {
		t.Fatal("expected zero-value defaults")
	}
}

func TestLoadParsesUpscaleSection(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.toml")
	os.WriteFile(p, []byte("[upscale]\ninput = \"~/Pictures/originals\"\nwidth = 4096\n"), 0o644)
	t.Setenv("FYNX_CONFIG", p)
	c, e := Load()
	if e != nil {
		t.Fatal(e)
	}
	if c.Upscale.Width != 4096 {
		t.Fatalf("width = %d", c.Upscale.Width)
	}
	if c.Upscale.Input != "~/Pictures/originals" {
		t.Fatalf("input = %q", c.Upscale.Input)
	}
}

func TestExpandResolvesTilde(t *testing.T) {
	h, _ := os.UserHomeDir()
	if got := Expand("~"); got != h {
		t.Fatalf("Expand(~) = %q", got)
	}
	if got := Expand("~/Pictures"); got != filepath.Join(h, "Pictures") {
		t.Fatalf("Expand(~/Pictures) = %q", got)
	}
	if got := Expand("/abs/path"); got != "/abs/path" {
		t.Fatalf("Absolute path modified: %q", got)
	}
}
