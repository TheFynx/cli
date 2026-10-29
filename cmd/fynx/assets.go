package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"fynx.local/cli/internal/config"
	"fynx.local/cli/internal/upscale"
)

func newAssetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Engine setup and wallpaper upscaling",
	}
	cmd.AddCommand(newAssetsSetupCmd(), newAssetsUpscaleCmd())
	return cmd
}

func newAssetsSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Install the pinned Real-ESRGAN engine and x4plus model",
		Long: "Downloads checksum-pinned official Linux x86_64 Real-ESRGAN NCNN/Vulkan\n" +
			"release binaries and the official x4plus model into the fynx data\n" +
			"directory. Requires curl, unzip, a working Vulkan loader, and\n" +
			"ImageMagick v7 (magick). Run ./build.sh first so the installer\n" +
			"script is in place.",
		RunE: func(*cobra.Command, []string) error {
			script := filepath.Join(config.DataDir(), "setup-realesrgan.sh")
			if _, e := os.Stat(script); e != nil {
				return errors.New("installer missing; run ./build.sh first")
			}
			return run("bash", script)
		},
	}
}

func newAssetsUpscaleCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "upscale",
		Short: "AI-upscale images 4x, then resize to a target width",
		Long: "Runs the installed Real-ESRGAN engine at 4x and exports at the target\n" +
			"width with Lanczos, preserving aspect ratio. Originals stay untouched;\n" +
			"each export carries a SHA-256 receipt that makes reruns cheap and\n" +
			"tamper detection automatic. Defaults may be set in config.toml.",
		RunE: func(c *cobra.Command, args []string) error {
			if e := loadConfig(); e != nil {
				return e
			}
			if len(args) != 0 {
				return errors.New("unexpected arguments; see fynx assets upscale --help")
			}
			p := mergeUpscale(c.Flags())
			if p.Input == "" || p.Output == "" {
				return errors.New("use upscale --input DIR --output DIR (or set [upscale] in config.toml)")
			}
			return upscale.Run(p)
		},
	}
	f := c.Flags()
	f.String("input", "", "source image directory")
	f.String("output", "", "separate export directory")
	f.String("engine", "", "Real-ESRGAN binary (default: fynx data dir)")
	f.String("model", "realesrgan-x4plus", "NCNN model name")
	f.String("models", "", "model directory (defaults beside engine)")
	f.Int("width", 3840, "export width; original aspect ratio is preserved")
	f.Int("tile", 128, "GPU tile size")
	f.Bool("dry-run", false, "show operations without creating outputs or running GPU inference")
	f.Bool("force", false, "regenerate matching cached exports")
	return c
}

// mergeUpscale resolves flag vs config-file vs built-in precedence.
func mergeUpscale(f *pflag.FlagSet) upscale.Params {
	p := upscale.Defaults()
	str := func(name string, cfg, def string) string {
		if f.Changed(name) {
			v, _ := f.GetString(name)
			return config.Expand(v)
		}
		if cfg != "" {
			return config.Expand(cfg)
		}
		return def
	}
	num := func(name string, cfg, def int) int {
		if f.Changed(name) {
			v, _ := f.GetInt(name)
			return v
		}
		if cfg != 0 {
			return cfg
		}
		return def
	}
	p.Input = str("input", config.C.Upscale.Input, "")
	p.Output = str("output", config.C.Upscale.Output, "")
	p.Engine = str("engine", config.C.Upscale.Engine, p.Engine)
	p.Models = str("models", config.C.Upscale.Models, "")
	p.Model = str("model", config.C.Upscale.Model, p.Model)
	p.Width = num("width", config.C.Upscale.Width, p.Width)
	p.Tile = num("tile", config.C.Upscale.Tile, p.Tile)
	p.Plan, _ = f.GetBool("dry-run")
	p.Force, _ = f.GetBool("force")
	return p
}
