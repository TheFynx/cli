package theme

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// installFixtures copies every testdata palette into a fake home so the
// multi-variant minecraft builder can run.
func installFixtures(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	for _, slug := range []string{"minecraft-deepstone", "minecraft-nether", "minecraft-end", "midgar-mako"} {
		b, e := os.ReadFile(filepath.Join("testdata", slug, "colors.toml"))
		if e != nil {
			t.Fatal(e)
		}
		if e := write(filepath.Join(home, ".config/omarchy/themes", slug, "colors.toml"), string(b)); e != nil {
			t.Fatal(e)
		}
	}
	return home
}

func TestBuildMinecraftGeneratesExpectedAssets(t *testing.T) {
	home := installFixtures(t)
	if e := Build(home, "minecraft"); e != nil {
		t.Fatal(e)
	}
	td := filepath.Join(home, ".config/omarchy/themes/minecraft-deepstone")
	css, e := os.ReadFile(filepath.Join(td, "gtk.css"))
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(string(css), "@define-color view_bg_color #101513;") {
		t.Fatal("gtk.css missing view roles")
	}
	qtc, e := os.ReadFile(filepath.Join(td, "qt.colors"))
	if e != nil {
		t.Fatal(e)
	}
	if !strings.HasPrefix(string(qtc), "[ColorScheme]\nactive_colors=#ffe2e7d7,") {
		t.Fatal("qt.colors wrong header or missing #ff prefix: " + strings.SplitN(string(qtc), "\n", 2)[1])
	}
	// The grass icon became the menu icon.
	m, e := os.ReadFile(filepath.Join(td, "assets/menu.svg"))
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(string(m), "#91c96a") {
		t.Fatal("menu.svg did not inherit the grass accent")
	}
	// Screensaver rows are centered within a 29-wide canvas (CPython quirk).
	ss, e := os.ReadFile(filepath.Join(td, "screensaver.txt"))
	if e != nil {
		t.Fatal(e)
	}
	for _, want := range []string{"       ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄", "M I N E C R A F T"} {
		if !strings.Contains(string(ss), want) {
			t.Fatalf("screensaver missing %q", want)
		}
	}
	// Icons were mapped into all three contexts.
	for _, p := range []string{
		"scalable/places/folder.svg", "scalable/apps/com.mitchellh.ghostty.svg", "scalable/status/dialog-error.svg",
	} {
		if _, e := os.Stat(filepath.Join(home, ".local/share/icons/Minecraft-Deepstone", p)); e != nil {
			t.Fatalf("missing icon %s: %v", p, e)
		}
	}
	// Unknown theme names fail cleanly.
	if e := Build(home, "nope"); e == nil {
		t.Fatal("unknown theme accepted")
	}
}

func TestBuildMidgarGeneratesExpectedAssets(t *testing.T) {
	home := installFixtures(t)
	if e := Build(home, "midgar"); e != nil {
		t.Fatal(e)
	}
	td := filepath.Join(home, ".config/omarchy/themes/midgar-mako")
	// folder.svg is icon-only, never written as an asset (as before).
	if _, e := os.Stat(filepath.Join(td, "assets/folder.svg")); e == nil {
		t.Fatal("folder.svg must not be written to assets")
	}
	menu, e := os.ReadFile(filepath.Join(td, "assets/menu.svg"))
	if e != nil {
		t.Fatal(e)
	}
	// The palette's accent happens to equal the authored glyph color, so
	// the authored bytes pass through unchanged.
	if !strings.Contains(string(menu), "#63e6be") {
		t.Fatalf("menu.svg lost the accent color: %v", e)
	}
	qtc, e := os.ReadFile(filepath.Join(td, "qt.colors"))
	if e != nil {
		t.Fatal(e)
	}
	if !strings.HasPrefix(string(qtc), "[ColorScheme]\nactive_colors=#ffd9e5e7,#ff19272e,") {
		t.Fatal("midgar qt.colors wrong: " + strings.SplitN(string(qtc), "\n", 2)[1])
	}
	// The empty status context dir still exists, matching the old builder.
	st := filepath.Join(home, ".local/share/icons/Midgar-Mako/scalable/status")
	if fi, e := os.Stat(st); e != nil || !fi.IsDir() {
		t.Fatal("status icon dir missing")
	}
}

func TestLoadRejectsIncompletePalette(t *testing.T) {
	home := t.TempDir()
	if e := write(filepath.Join(home, ".config/omarchy/themes/midgar-mako", "colors.toml"), "accent = \"#000\"\n"); e != nil {
		t.Fatal(e)
	}
	if e := Build(home, "midgar"); e == nil {
		t.Fatal("loading succeeded with no colors")
	}
}
