# Fynx CLI

Personal Go CLI that replaced `~/.functions` with maintained subcommands.
Source lives at https://github.com/TheFynx/cli. Clone it to
`~/git/thefynx/cli`, which is where provisioning expects it. The
executable is `fynx`. Theme builders, disk reports, DNS lookups,
reachability probes, and repo tools run natively in Go. They need no
Python, no Pillow, and no dig. Engine upscaling still shells to the
official Real-ESRGAN binary.

```bash
go test ./...
./build.sh
fynx --help
```

`build.sh` compiles into `~/.local/bin/fynx` and installs the persistent
native asset installer under `${XDG_DATA_HOME:-~/.local/share}/fynx`.
Keep `~/.local/bin` on PATH. The build does not remove old shell
wrappers. Retire them yourself: run `fynx aliases bash` (or `zsh`,
`fish`), paste the printed block into your shell rc, then delete
`~/.functions`.

## Commands

| Command | Replaces | Notes |
|---|---|---|
| `fynx assets setup` | new | checksum-pinned Real-ESRGAN + x4plus install |
| `fynx assets upscale` | new | receipt-verified 4x upscaling and export |
| `fynx theme build minecraft\|midgar` | two Python builders | native SVG/GTK/Qt generation |
| `fynx docker clean [all\|fire]` | `dclean` | `--dry-run` prints planned commands |
| `fynx disk usage [PATH]...` | `fs` | `du -sbh`-style apparent sizes, no subprocess |
| `fynx dns lookup NAME` | `digga` | built-in resolver, grouped records, no dig |
| `fynx net isup URI` | `isup` | real HTTP probe, exits 1 when down |
| `fynx man [SECTION] PAGE...` | colored `man` | alias `man='fynx man'` |
| `fynx repo check BASE` | `git-check` | per-repo report plus summary, no symlink loops |
| `fynx repo contents DIR` | `copy_repo_contents` | inlines cmd/ and internal/ file contents |
| `fynx aliases [shell]` | new | prints the alias block above |

Two shell bugs are corrected by design. `isup` used to report the
opposite verdict for HTTP 200. The `copy_repo_contents` tar step
silently re-extracted `cmd/` and `internal/` under `/` instead of
appending their contents.

## Configuration

Configuration file discovery order: `$FYNX_CONFIG`, then
`$XDG_CONFIG_HOME/fynx/config.toml`, then `~/.config/fynx/config.toml`.
A missing file means built-in defaults. Flags override config, and
config overrides defaults. Values accept `~`.

```toml
[upscale]
input  = "~/Pictures/originals"
output = "~/Pictures/upscaled"
# engine = "~/data/realesrgan/realesrgan-ncnn-vulkan"
# models = "~/data/realesrgan/models"
# model  = "realesrgan-x4plus"
width  = 3840
tile   = 128
```

## Wallpaper upscaling

```bash
fynx assets setup
fynx assets upscale --dry-run
fynx assets upscale
```

Setup installs the checksum-pinned official Real-ESRGAN NCNN/Vulkan
release binaries for Linux x86_64 and the official x4plus model. It
requires curl, unzip, a working Vulkan driver/loader, and ImageMagick
v7 (`magick`). The recorded hashes pin downloaded bytes. They are not
publisher signatures. Select a custom engine or model directory with
`--engine` and `--models`, or in config.

Upscaling runs x4plus at 4x, then a Lanczos resize to 3840 pixels wide
by default. The original aspect ratio is preserved with nearest-integer
export height. `--width` changes the target. Files export separately.
Originals remain unchanged. Intermediates are removed after success.
One SHA-256 receipt per export records the source, engine/model bytes,
dimensions, and output. A matching receipt skips repeated inference.
Changed originals, changed models, or tampered exports regenerate.
`--force` bypasses the cache. Outputs with different source extensions
get distinct names. Work runs serially with bounded GPU tiles.

## Native theme builders

```bash
fynx theme build midgar
fynx theme build minecraft
```

Generates native SVG icons, palette-driven GTK/Qt assets, and
screensaver text from the installed themes' `colors.toml`. The output is
byte-compatible with the retired Python builders, verified by tree
diff. The commands overwrite generated assets in the installed theme.
Review handcrafted edits before you run them, and reapply with
`omarchy theme set <slug>` afterward. Runtime theme hooks remain in the
blueprints' common theme layer. No wallpapers, no image provider, no
network.

Omarchy's packaged CLI scans its own binary directory. Its shell plugins
provide UI, menu, and lock integration. `fynx` remains an independent
general CLI and does not modify `/usr/share/omarchy`. Cursor appearance
is left to normal desktop settings.