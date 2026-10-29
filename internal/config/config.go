// Package config loads fynx settings from an optional TOML file.
//
// File discovery order: $FYNX_CONFIG, then $XDG_CONFIG_HOME/fynx/config.toml,
// then ~/.config/fynx/config.toml. A missing file is not an error; defaults
// apply. Values may reference the home directory with a leading ~.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Upscale holds defaults for `fynx assets upscale`. Command-line flags
// override these; these override built-in defaults.
type Upscale struct {
	Input  string // source image directory
	Output string // export directory
	Engine string // Real-ESRGAN executable path
	Models string // model directory (defaults beside the engine)
	Model  string // NCNN model name
	Width  int    // export width in pixels
	Tile   int    // GPU tile size
}

// Config is the fynx configuration file schema.
type Config struct {
	Upscale Upscale `toml:"upscale"`
}

// C is the process-wide configuration, populated by Load at startup.
var C Config

// Path returns where the configuration file is expected.
func Path() string {
	if x := os.Getenv("FYNX_CONFIG"); x != "" {
		return x
	}
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		h, _ := os.UserHomeDir()
		dir = filepath.Join(h, ".config")
	}
	return filepath.Join(dir, "fynx", "config.toml")
}

// Load reads the configuration file, returning defaults when absent.
func Load() (Config, error) {
	p := Path()
	b, err := os.ReadFile(p)
	if errors.Is(err, fs.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}
	var c Config
	if e := toml.Unmarshal(b, &c); e != nil {
		return Config{}, fmt.Errorf("%s: %w", p, e)
	}
	return c, nil
}

// DataDir is the persistent state directory: engine binaries, downloads.
func DataDir() string {
	if x := os.Getenv("XDG_DATA_HOME"); x != "" {
		return filepath.Join(x, "fynx")
	}
	h, _ := os.UserHomeDir()
	return filepath.Join(h, ".local/share/fynx")
}

// Expand resolves a leading ~ in a configuration path against the home
// directory. Absolute and relative paths pass through unchanged.
func Expand(p string) string {
	if p == "~" || filepath.Base(p) == p && len(p) > 0 && p[0] == '~' && len(p) == 1 {
		// p == "~"
	}
	if p == "~" {
		h, _ := os.UserHomeDir()
		return h
	}
	if len(p) >= 2 && p[0] == '~' && (p[1] == '/' || p[1] == filepath.Separator) {
		h, _ := os.UserHomeDir()
		return filepath.Join(h, p[2:])
	}
	return p
}
