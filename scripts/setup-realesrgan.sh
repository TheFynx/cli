#!/usr/bin/env bash
# Pinned official xinntao release assets. Digests are recorded from those downloads.
set -euo pipefail
[ "$(uname -s)/$(uname -m)" = "Linux/x86_64" ] || { echo "Pinned engine requires Linux x86_64" >&2; exit 1; }
root="${XDG_DATA_HOME:-$HOME/.local/share}/fynx/realesrgan"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
curl -fL --retry 2 'https://github.com/xinntao/Real-ESRGAN-ncnn-vulkan/releases/download/v0.2.0/realesrgan-ncnn-vulkan-v0.2.0-ubuntu.zip' -o "$work/engine.zip"
curl -fL --retry 2 'https://github.com/xinntao/Real-ESRGAN/releases/download/v0.2.5.0/realesrgan-ncnn-vulkan-20220424-ubuntu.zip' -o "$work/models.zip"
printf '%s  %s\n' d0e8e1cf954f5cde11be4745dd912cc3774bef36f71c5b1cb8f74c4112b6e919 "$work/engine.zip" e5aa6eb131234b87c0c51f82b89390f5e3e642b7b70f2b9bbe95b6a285a40c96 "$work/models.zip" | sha256sum -c -
unzip -q "$work/engine.zip" -d "$work/engine"
unzip -q "$work/models.zip" -d "$work/models"
install -d "$root/models"
engine=$(find "$work/engine" -type f -name realesrgan-ncnn-vulkan -print -quit)
install -m 755 "$engine" "$root/realesrgan-ncnn-vulkan"
for name in realesrgan-x4plus.bin realesrgan-x4plus.param; do
 file=$(find "$work/models" -type f -name "$name" -print -quit)
 install -m 644 "$file" "$root/models/$name"
done
find "$work/engine" "$work/models" -type f -iname '*license*' -exec cp -t "$root" {} +
echo "Installed official engine and x4plus model to $root"
