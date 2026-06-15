#!/usr/bin/env bash
# Build admin Tailwind CSS into web/static/css/tailwind.css (embedded at go build).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

OS="$(uname -s)"
ARCH="$(uname -m)"
case "$OS-$ARCH" in
	Darwin-arm64) TW_BIN="tailwindcss-macos-arm64" ;;
	Darwin-x86_64) TW_BIN="tailwindcss-macos-x64" ;;
	Linux-x86_64) TW_BIN="tailwindcss-linux-x64" ;;
	Linux-aarch64) TW_BIN="tailwindcss-linux-arm64" ;;
	*)
		echo "Unsupported platform: $OS $ARCH"
		exit 1
		;;
esac

TW="/tmp/tailwindcss-$$"
curl -fsSL "https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/${TW_BIN}" -o "$TW"
chmod +x "$TW"
"$TW" -i web/static/css/admin-input.css -o web/static/css/tailwind.css --minify
rm -f "$TW"
echo "Built web/static/css/tailwind.css"
