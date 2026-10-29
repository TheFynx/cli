// Package theme builds native omarchy theme assets from the installed
// themes' colors.toml. Generation is fully offline and deterministic:
// icons are authored SVG strings, GTK/Qt rules are color substitution.
package theme

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Colors mirrors omarchy theme colors.toml keys.
type Colors struct {
	Mode              string `toml:"mode"`
	Accent            string `toml:"accent"`
	Selection         string `toml:"selection"`
	Muted             string `toml:"muted"`
	Background        string `toml:"background"`
	DarkBackground    string `toml:"dark_background"`
	DarkerBackground  string `toml:"darker_background"`
	LighterBackground string `toml:"lighter_background"`
	Foreground        string `toml:"foreground"`
	DarkForeground    string `toml:"dark_foreground"`
	LightForeground   string `toml:"light_foreground"`
	BrightForeground  string `toml:"bright_foreground"`
	Red               string `toml:"red"`
	Yellow            string `toml:"yellow"`
	Orange            string `toml:"orange"`
	Green             string `toml:"green"`
	Cyan              string `toml:"cyan"`
	Blue              string `toml:"blue"`
	Magenta           string `toml:"magenta"`
	Brown             string `toml:"brown"`
	BrightRed         string `toml:"bright_red"`
	BrightYellow      string `toml:"bright_yellow"`
	BrightGreen       string `toml:"bright_green"`
	BrightCyan        string `toml:"bright_cyan"`
	BrightBlue        string `toml:"bright_blue"`
	BrightMagenta     string `toml:"bright_magenta"`
}

// Build regenerates assets for a named theme family under home.
func Build(home, name string) error {
	switch name {
	case "minecraft":
		return buildMinecraft(home)
	case "midgar":
		return buildMidgar(home)
	}
	return fmt.Errorf("unknown theme %q (use minecraft or midgar)", name)
}

// load reads ~/.config/omarchy/themes/<slug>/colors.toml.
func load(home, slug string) (Colors, error) {
	p := filepath.Join(home, ".config/omarchy/themes", slug, "colors.toml")
	b, err := os.ReadFile(p)
	if err != nil {
		return Colors{}, fmt.Errorf("theme %s: %w", slug, err)
	}
	var c Colors
	if e := toml.Unmarshal(b, &c); e != nil {
		return Colors{}, fmt.Errorf("%s: %w", p, e)
	}
	if c.Accent == "" || c.Background == "" || c.Foreground == "" {
		return Colors{}, fmt.Errorf("%s: accent/background/foreground required", p)
	}
	return c, nil
}

// colorsPairs exposes every non-empty value as key/value pairs for
// placeholder substitution in templates such as Qt stylesheets.
func (c Colors) pairs() []struct{ key, value string } {
	type pair = struct{ key, value string }
	all := []pair{
		{"mode", c.Mode}, {"accent", c.Accent}, {"selection", c.Selection},
		{"muted", c.Muted}, {"background", c.Background},
		{"dark_background", c.DarkBackground}, {"darker_background", c.DarkerBackground},
		{"lighter_background", c.LighterBackground}, {"foreground", c.Foreground},
		{"dark_foreground", c.DarkForeground}, {"light_foreground", c.LightForeground},
		{"bright_foreground", c.BrightForeground}, {"red", c.Red}, {"yellow", c.Yellow},
		{"orange", c.Orange}, {"green", c.Green}, {"cyan", c.Cyan}, {"blue", c.Blue},
		{"magenta", c.Magenta}, {"brown", c.Brown}, {"bright_red", c.BrightRed},
		{"bright_yellow", c.BrightYellow}, {"bright_green", c.BrightGreen},
		{"bright_cyan", c.BrightCyan}, {"bright_blue", c.BrightBlue},
		{"bright_magenta", c.BrightMagenta},
	}
	kept := all[:0]
	for _, p := range all {
		if p.value != "" {
			kept = append(kept, p)
		}
	}
	return kept
}

func write(path, content string) error {
	if e := os.MkdirAll(filepath.Dir(path), 0o755); e != nil {
		return e
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// title capitalizes the first letter and lowers the rest, like Python's
// str.title on a single word: "deepstone" -> "Deepstone".
func title(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func copyFile(src, dst string) error {
	b, e := os.ReadFile(src)
	if e != nil {
		return e
	}
	return write(dst, string(b))
}
