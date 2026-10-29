package upscale

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestUpscaleOriginalsAndVerifiedCache(t *testing.T) {
	if _, e := exec.LookPath("magick"); e != nil {
		t.Skip("ImageMagick unavailable")
	}
	root := t.TempDir()
	input := filepath.Join(root, "in")
	output := filepath.Join(root, "out")
	models := filepath.Join(root, "models")
	os.MkdirAll(input, 0o755)
	os.MkdirAll(models, 0o755)
	source := filepath.Join(input, "scene.png")
	f, _ := os.Create(source)
	im := image.NewRGBA(image.Rect(0, 0, 8, 5))
	im.Set(1, 1, color.RGBA{255, 0, 0, 255})
	png.Encode(f, im)
	f.Close()
	before, _ := hash(source)
	for _, name := range []string{"realesrgan-x4plus.bin", "realesrgan-x4plus.param"} {
		os.WriteFile(filepath.Join(models, name), []byte("fixture"), 0o644)
	}
	engine := filepath.Join(root, "engine")
	calls := filepath.Join(root, "calls")
	script := "#!/bin/sh\nwhile [ $# -gt 0 ]; do case $1 in -i) shift; input=$1;; -o) shift; output=$1;; esac; shift; done\nprintf x >> '" + calls + "'\ncp \"$input\" \"$output\"\n"
	os.WriteFile(engine, []byte(script), 0o755)
	p := Params{Input: input, Output: output, Engine: engine, Models: models, Width: 32, Model: "realesrgan-x4plus"}
	if e := Run(p); e != nil {
		t.Fatal(e)
	}
	w, h, e := dimensions(filepath.Join(output, "scene.png"))
	if e != nil || w != 32 || h != 20 {
		t.Fatalf("wrong export: %dx%d %v", w, h, e)
	}
	after, _ := hash(source)
	if before != after {
		t.Fatal("original changed")
	}
	// A valid cache must not re-run inference.
	if e := Run(p); e != nil {
		t.Fatal(e)
	}
	if b, _ := os.ReadFile(calls); len(b) != 1 {
		t.Fatal("valid cache repeated inference")
	}
	// A tampered export must be regenerated.
	os.WriteFile(filepath.Join(output, "scene.png"), []byte("tampered"), 0o644)
	if e := Run(p); e != nil {
		t.Fatal(e)
	}
	if b, _ := os.ReadFile(calls); len(b) != 2 {
		t.Fatal("tampered export was trusted")
	}
	// Refuse same source and output.
	if e := Run(Params{Input: input, Output: input, Engine: engine, Models: models, Width: 32}); e == nil {
		t.Fatal("same source/output accepted")
	}
	// Dry run writes nothing.
	plan := filepath.Join(root, "plan")
	if e := Run(Params{Input: input, Output: plan, Plan: true, Engine: engine, Models: models, Width: 32}); e != nil {
		t.Fatal(e)
	}
	if _, e := os.Stat(plan); !os.IsNotExist(e) {
		t.Fatal("dry run wrote output")
	}
}
