package main

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
	os.MkdirAll(input, 0755)
	os.MkdirAll(models, 0755)
	source := filepath.Join(input, "scene.png")
	f, _ := os.Create(source)
	im := image.NewRGBA(image.Rect(0, 0, 8, 5))
	im.Set(1, 1, color.RGBA{255, 0, 0, 255})
	png.Encode(f, im)
	f.Close()
	before, _ := hash(source)
	for _, name := range []string{"realesrgan-x4plus.bin", "realesrgan-x4plus.param"} {
		os.WriteFile(filepath.Join(models, name), []byte("fixture"), 0644)
	}
	engine := filepath.Join(root, "engine")
	calls := filepath.Join(root, "calls")
	script := "#!/bin/sh\nwhile [ $# -gt 0 ]; do case $1 in -i) shift; input=$1;; -o) shift; output=$1;; esac; shift; done\nprintf x >> '" + calls + "'\ncp \"$input\" \"$output\"\n"
	os.WriteFile(engine, []byte(script), 0755)
	args := []string{"--input", input, "--output", output, "--engine", engine, "--models", models, "--width", "32"}
	if e := upscale(args); e != nil {
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
	if e := upscale(args); e != nil {
		t.Fatal(e)
	}
	b, _ := os.ReadFile(calls)
	if len(b) != 1 {
		t.Fatal("valid cache repeated inference")
	}
	os.WriteFile(filepath.Join(output, "scene.png"), []byte("tampered"), 0644)
	if e := upscale(args); e != nil {
		t.Fatal(e)
	}
	b, _ = os.ReadFile(calls)
	if len(b) != 2 {
		t.Fatal("tampered export was trusted")
	}
	if e := upscale([]string{"--input", input, "--output", input}); e == nil {
		t.Fatal("same source/output accepted")
	}
	plan := filepath.Join(root, "plan")
	if e := upscale([]string{"--input", input, "--output", plan, "--dry-run"}); e != nil {
		t.Fatal(e)
	}
	if _, e := os.Stat(plan); !os.IsNotExist(e) {
		t.Fatal("dry run wrote output")
	}
}
