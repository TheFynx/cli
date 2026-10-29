// Package upscale runs the Real-ESRGAN engine and exports resized images
// with verified SHA-256 receipts. GPU inference stays in the official
// engine binary; this package only orchestrates it.
package upscale

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/TheFynx/cli/internal/config"
)

// Receipt records exactly what produced an export, so a matching receipt
// proves the export is current.
type Receipt struct {
	SourceSHA string `json:"source_sha256"`
	Model     string `json:"model"`
	EngineSHA string `json:"engine_sha256"`
	ModelSHA  string `json:"model_sha256"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	OutputSHA string `json:"output_sha256"`
}

// Params carries one upscale invocation. Zero values fall back to the
// built-in defaults; the caller merges config-file values first.
type Params struct {
	Input  string
	Output string
	Engine string
	Models string
	Model  string
	Width  int
	Tile   int
	Plan   bool // dry run
	Force  bool
}

func hash(path string) (string, error) {
	f, e := os.Open(path)
	if e != nil {
		return "", e
	}
	defer f.Close()
	h := sha256.New()
	if _, e = io.Copy(h, f); e != nil {
		return "", e
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func dimensions(path string) (int, int, error) {
	b, e := exec.Command("magick", "identify", "-ping", "-format", "%w %h", path).Output()
	if e != nil {
		return 0, 0, e
	}
	parts := strings.Fields(string(b))
	if len(parts) != 2 {
		return 0, 0, errors.New("expected a single image")
	}
	w, e := strconv.Atoi(parts[0])
	if e != nil {
		return 0, 0, e
	}
	h, e := strconv.Atoi(parts[1])
	if e != nil || w <= 0 || h <= 0 {
		return 0, 0, errors.New("invalid image dimensions")
	}
	return w, h, nil
}

func canonical(path string) (string, error) {
	a, e := filepath.Abs(path)
	if e != nil {
		return "", e
	}
	if resolved, e := filepath.EvalSymlinks(a); e == nil {
		return resolved, nil
	}
	parent, e := filepath.EvalSymlinks(filepath.Dir(a))
	if e != nil {
		return a, nil
	}
	return filepath.Join(parent, filepath.Base(a)), nil
}

// Run upscales every supported image in Params.Input into Params.Output.
func Run(p Params) error {
	if p.Input == "" || p.Output == "" {
		return errors.New("use upscale --input DIR --output DIR [--width 3840]")
	}
	if p.Width <= 0 {
		p.Width = 3840
	}
	if p.Tile <= 0 {
		p.Tile = 128
	}
	input, e := canonical(p.Input)
	if e != nil {
		return e
	}
	output, e := canonical(p.Output)
	if e != nil {
		return e
	}
	if input == output {
		return errors.New("input and output must differ; original files are never overwritten")
	}
	entries, e := os.ReadDir(input)
	if e != nil {
		return e
	}
	exe, e := filepath.Abs(p.Engine)
	if e != nil {
		return e
	}
	modelDir := p.Models
	if modelDir == "" {
		modelDir = filepath.Join(filepath.Dir(exe), "models")
	}
	modelDir, e = filepath.Abs(modelDir)
	if e != nil {
		return e
	}
	var engineHash, modelHash string
	if !p.Plan {
		if _, e = exec.LookPath("magick"); e != nil {
			return e
		}
		engineHash, e = hash(exe)
		if e != nil {
			return fmt.Errorf("engine unavailable; run fynx assets setup: %w", e)
		}
		a, e := hash(filepath.Join(modelDir, p.Model+".bin"))
		if e != nil {
			return e
		}
		b, e := hash(filepath.Join(modelDir, p.Model+".param"))
		if e != nil {
			return e
		}
		modelHash = a + b
		if e = os.MkdirAll(output, 0o755); e != nil {
			return e
		}
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" {
			continue
		}
		count++
		source := filepath.Join(input, entry.Name())
		target := filepath.Join(output, entry.Name()+".png")
		if ext == ".png" {
			target = filepath.Join(output, entry.Name())
		}
		if source == target {
			return errors.New("refusing to overwrite source")
		}
		w, h, e := dimensions(source)
		if e != nil {
			return fmt.Errorf("%s: %w", source, e)
		}
		height := int(float64(h)*float64(p.Width)/float64(w) + 0.5)
		fmt.Printf("%s: %dx%d -> %dx%d (%s)\n", entry.Name(), w, h, p.Width, height, p.Model)
		if p.Plan {
			continue
		}
		sh, e := hash(source)
		if e != nil {
			return e
		}
		r := Receipt{SourceSHA: sh, Model: p.Model, EngineSHA: engineHash, ModelSHA: modelHash, Width: p.Width, Height: height}
		receipt := target + ".json"
		var old Receipt
		b, _ := os.ReadFile(receipt)
		_ = json.Unmarshal(b, &old)
		oh, _ := hash(target)
		if !p.Force && old.SourceSHA == r.SourceSHA && old.Model == r.Model &&
			old.EngineSHA == r.EngineSHA && old.ModelSHA == r.ModelSHA &&
			old.Width == r.Width && old.Height == r.Height &&
			old.OutputSHA != "" && old.OutputSHA == oh {
			fmt.Println("  cached, verified")
			continue
		}
		work, e := os.MkdirTemp(output, ".upscale-")
		if e != nil {
			return e
		}
		defer os.RemoveAll(work)
		native := filepath.Join(work, "inference.png")
		export := filepath.Join(work, "export.png")
		c := exec.Command(exe, "-i", source, "-o", native, "-m", modelDir, "-n", p.Model, "-s", "4", "-t", strconv.Itoa(p.Tile), "-j", "1:1:1")
		c.Dir = filepath.Dir(exe)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if e = c.Run(); e != nil {
			return e
		}
		if e = run("magick", native, "-filter", "Lanczos", "-resize", fmt.Sprintf("%dx%d!", p.Width, height), export); e != nil {
			return e
		}
		ew, eh, e := dimensions(export)
		if e != nil || ew != p.Width || eh != height {
			return errors.New("export dimensions failed verification")
		}
		r.OutputSHA, e = hash(export)
		if e != nil {
			return e
		}
		if e = os.Rename(export, target); e != nil {
			return e
		}
		b, e = json.MarshalIndent(r, "", "  ")
		if e != nil {
			return e
		}
		if e = os.WriteFile(receipt, append(b, '\n'), 0o644); e != nil {
			return e
		}
	}
	if count == 0 {
		return errors.New("no supported images found")
	}
	return nil
}

func run(name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// Defaults returns built-in fallbacks; used when neither the config file
// nor flags supplied a value.
func Defaults() Params {
	return Params{
		Engine: filepath.Join(config.DataDir(), "realesrgan/realesrgan-ncnn-vulkan"),
		Model:  "realesrgan-x4plus",
		Width:  3840,
		Tile:   128,
	}
}
