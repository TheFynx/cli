#!/usr/bin/env bash
set -euo pipefail
root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
install -d "$HOME/.local/bin" "${XDG_DATA_HOME:-$HOME/.local/share}/fynx"
cd "$root"
GOTOOLCHAIN=local go build -buildvcs=false -trimpath -o "$HOME/.local/bin/fynx.new" ./cmd/fynx
mv "$HOME/.local/bin/fynx.new" "$HOME/.local/bin/fynx"
install -m 644 scripts/setup-realesrgan.sh "${XDG_DATA_HOME:-$HOME/.local/share}/fynx/"
# Theme builders are compiled in now; retire loose scripts from old installs.
rm -rf "${XDG_DATA_HOME:-$HOME/.local/share}/fynx/assets"
rm -f "${XDG_DATA_HOME:-$HOME/.local/share}/fynx/build-midgar.py" \
      "${XDG_DATA_HOME:-$HOME/.local/share}/fynx/build-minecraft.py"