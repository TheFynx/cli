#!/usr/bin/env bash
set -euo pipefail
root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
install -d "$HOME/.local/bin" "${XDG_DATA_HOME:-$HOME/.local/share}/fynx"
cd "$root"
GOTOOLCHAIN=local go build -buildvcs=false -trimpath -o "$HOME/.local/bin/fynx.new" ./cmd/fynx
mv "$HOME/.local/bin/fynx.new" "$HOME/.local/bin/fynx"
install -m 644 scripts/build-*.py scripts/setup-realesrgan.sh "${XDG_DATA_HOME:-$HOME/.local/share}/fynx/"

install -d "${XDG_DATA_HOME:-$HOME/.local/share}/fynx/assets"
install -m 644 scripts/assets/*.svg "${XDG_DATA_HOME:-$HOME/.local/share}/fynx/assets/"
