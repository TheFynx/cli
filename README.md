# Fynx CLI

Personal Go CLI for tools that belong in one maintained program rather than
shell/dotfile functions. The checkout is `~/git/thefynx/cli`; the executable is
`fynx`. No remote repository URL is assumed yet. The Go CLI has no third-party
module dependencies and works without a desktop for asset processing.

```bash
go test ./cmd/fynx
./build.sh
fynx --help
```

`build.sh` compiles into `~/.local/bin/fynx` and installs the persistent native
asset builders under `${XDG_DATA_HOME:-~/.local/share}/fynx`. Ensure
`~/.local/bin` is on PATH. Existing shell functions are not removed implicitly;
migrate each behavior into a reviewed command, then retire its old wrapper.

## Wallpaper upscaling

```bash
fynx assets setup
fynx assets upscale --input ~/Pictures/originals --output ~/Pictures/upscaled --dry-run
fynx assets upscale --input ~/Pictures/originals --output ~/Pictures/upscaled
```

Setup installs checksum-pinned official Linux x86_64 Real-ESRGAN NCNN/Vulkan
release binaries and the official x4plus model. Requires curl, unzip, a working
Vulkan driver/loader, and ImageMagick. The recorded hashes pin downloaded bytes;
they are not publisher signatures. A custom engine/model directory can be
selected with `--engine` and `--models`.

Upscaling uses x4plus at 4× followed by Lanczos to 3840 pixels wide by default.
The original aspect ratio is preserved, with nearest-integer export height;
`--width` changes that target. Files are exported separately, originals remain
unchanged, and full intermediate images are removed after success. This is AI
upscaling, not native-resolution generation. One SHA-256 receipt per export
records the source, engine/model bytes, dimensions and final output. A matching
receipt skips repeated inference; changed originals/models or tampered exports
are regenerated. `--force` bypasses that cache. Outputs with different source
extensions receive distinct names. Work runs serially with bounded GPU tiles.

## Native theme builders

```bash
fynx theme build midgar
fynx theme build minecraft
```

Persistent Python builders generate native SVG/icons, palette-driven GTK/Qt assets from the installed theme colors.
Requires Python 3.11+ and Pillow. They do not generate wallpapers,
use an image provider, or call external APIs. These commands overwrite generated
assets in the installed theme; review handcrafted edits before running them,
and reapply with `omarchy theme set <slug>` afterward. Runtime theme hooks remain
in the blueprints' common theme layer.

Omarchy's packaged CLI scans its own binary directory; its shell plugins provide
UI/menu/lock integration. `fynx` remains an independent general CLI and does not
modify `/usr/share/omarchy`.

Cursor appearance is left to the normal desktop settings. Theme hooks and
asset builders do not install or select custom cursors.
