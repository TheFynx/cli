// fynx is a local asset CLI; GPU inference stays in the official Real-ESRGAN engine.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Receipt struct {
	SourceSHA string `json:"source_sha256"`
	Model     string `json:"model"`
	EngineSHA string `json:"engine_sha256"`
	ModelSHA  string `json:"model_sha256"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	OutputSHA string `json:"output_sha256"`
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
func command(name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
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
func dataDir() string {
	h, _ := os.UserHomeDir()
	if x := os.Getenv("XDG_DATA_HOME"); x != "" {
		return filepath.Join(x, "fynx")
	}
	return filepath.Join(h, ".local/share/fynx")
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
func upscale(args []string) error {
	fs := flag.NewFlagSet("upscale", flag.ContinueOnError)
	in := fs.String("input", "", "source image directory")
	out := fs.String("output", "", "separate export directory")
	engine := fs.String("engine", filepath.Join(dataDir(), "realesrgan/realesrgan-ncnn-vulkan"), "Real-ESRGAN binary")
	model := fs.String("model", "realesrgan-x4plus", "NCNN model name")
	models := fs.String("models", "", "model directory (defaults beside engine)")
	width := fs.Int("width", 3840, "export width; original aspect ratio is preserved")
	tile := fs.Int("tile", 128, "GPU tile size")
	plan := fs.Bool("dry-run", false, "show operations without creating outputs or running GPU inference")
	force := fs.Bool("force", false, "regenerate matching cached exports")
	if e := fs.Parse(args); e != nil {
		return e
	}
	if *in == "" || *out == "" || *width <= 0 || *tile <= 0 || fs.NArg() != 0 {
		return errors.New("use upscale --input DIR --output DIR [--width 3840]")
	}
	input, e := canonical(*in)
	if e != nil {
		return e
	}
	output, e := canonical(*out)
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
	exe, e := filepath.Abs(*engine)
	if e != nil {
		return e
	}
	if *models == "" {
		*models = filepath.Join(filepath.Dir(exe), "models")
	}
	modelDir, e := filepath.Abs(*models)
	if e != nil {
		return e
	}
	var engineHash, modelHash string
	if !*plan {
		if _, e = exec.LookPath("magick"); e != nil {
			return e
		}
		engineHash, e = hash(exe)
		if e != nil {
			return fmt.Errorf("engine unavailable; run fynx assets setup: %w", e)
		}
		a, e := hash(filepath.Join(modelDir, *model+".bin"))
		if e != nil {
			return e
		}
		b, e := hash(filepath.Join(modelDir, *model+".param"))
		if e != nil {
			return e
		}
		modelHash = a + b
		if e = os.MkdirAll(output, 0755); e != nil {
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
		height := int(float64(h)*float64(*width)/float64(w) + 0.5)
		fmt.Printf("%s: %dx%d -> %dx%d (%s)\n", entry.Name(), w, h, *width, height, *model)
		if *plan {
			continue
		}
		sh, e := hash(source)
		if e != nil {
			return e
		}
		r := Receipt{SourceSHA: sh, Model: *model, EngineSHA: engineHash, ModelSHA: modelHash, Width: *width, Height: height}
		receipt := target + ".json"
		var old Receipt
		b, _ := os.ReadFile(receipt)
		_ = json.Unmarshal(b, &old)
		oh, _ := hash(target)
		if !*force && old.SourceSHA == r.SourceSHA && old.Model == r.Model && old.EngineSHA == r.EngineSHA && old.ModelSHA == r.ModelSHA && old.Width == r.Width && old.Height == r.Height && old.OutputSHA != "" && old.OutputSHA == oh {
			fmt.Println("  cached, verified")
			continue
		}
		work, e := os.MkdirTemp(output, ".upscale-")
		if e != nil {
			return e
		}
		native := filepath.Join(work, "inference.png")
		export := filepath.Join(work, "export.png")
		c := exec.Command(exe, "-i", source, "-o", native, "-m", modelDir, "-n", *model, "-s", "4", "-t", strconv.Itoa(*tile), "-j", "1:1:1")
		c.Dir = filepath.Dir(exe)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if e = c.Run(); e != nil {
			os.RemoveAll(work)
			return e
		}
		if e = command("magick", native, "-filter", "Lanczos", "-resize", fmt.Sprintf("%dx%d!", *width, height), export); e != nil {
			os.RemoveAll(work)
			return e
		}
		ew, eh, e := dimensions(export)
		if e != nil || ew != *width || eh != height {
			os.RemoveAll(work)
			return errors.New("export dimensions failed verification")
		}
		r.OutputSHA, e = hash(export)
		if e != nil {
			os.RemoveAll(work)
			return e
		}
		if e = os.Rename(export, target); e != nil {
			os.RemoveAll(work)
			return e
		}
		b, e = json.MarshalIndent(r, "", "  ")
		if e != nil {
			return e
		}
		if e = os.WriteFile(receipt, append(b, '\n'), 0644); e != nil {
			return e
		}
		os.RemoveAll(work)
	}
	if count == 0 {
		return errors.New("no supported images found")
	}
	return nil
}
func main() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Println("fynx assets upscale --input DIR --output DIR [--width 3840] [--dry-run]\n  Keeps originals unchanged; verifies cached exports using SHA-256 receipts.\nfynx assets setup  # install pinned official Real-ESRGAN + x4plus model\nfynx theme build minecraft|midgar  # regenerate installed assets; review manual edits first")
		return
	}
	var e error
	switch {
	case len(args) >= 2 && args[0] == "assets" && args[1] == "upscale":
		e = upscale(args[2:])
	case len(args) == 2 && args[0] == "assets" && args[1] == "setup":
		e = command("bash", filepath.Join(dataDir(), "setup-realesrgan.sh"))
	case len(args) == 3 && args[0] == "theme" && args[1] == "build" && (args[2] == "minecraft" || args[2] == "midgar"):
		e = command("python3", filepath.Join(dataDir(), "build-"+args[2]+".py"))
	default:
		e = errors.New("unknown command; run fynx --help")
	}
	if e != nil {
		fmt.Fprintln(os.Stderr, "fynx:", e)
		os.Exit(1)
	}
}
