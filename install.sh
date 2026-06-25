#!/usr/bin/env sh
set -eu

repo="${SEEKAI_REPO:-magalab/seekai-cli}"
version="${SEEKAI_VERSION:-latest}"
bin_dir="${SEEKAI_INSTALL_DIR:-/usr/local/bin}"
binary="seekai"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

case "$os" in
  darwin) goos="darwin" ;;
  linux) goos="linux" ;;
  *) echo "unsupported OS: $os" >&2; exit 1 ;;
esac

case "$arch" in
  x86_64|amd64) goarch="amd64" ;;
  arm64|aarch64) goarch="arm64" ;;
  *) echo "unsupported architecture: $arch" >&2; exit 1 ;;
esac

asset="seekai_${goos}_${goarch}"

if [ "$version" = "latest" ]; then
  url="https://github.com/${repo}/releases/latest/download/${asset}"
else
  url="https://github.com/${repo}/releases/download/${version}/${asset}"
fi

tmp="${TMPDIR:-/tmp}/${asset}.$$"
cleanup() {
  rm -f "$tmp"
}
trap cleanup EXIT

echo "downloading ${url}"
if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$url" -o "$tmp"
elif command -v wget >/dev/null 2>&1; then
  wget -qO "$tmp" "$url"
else
  echo "curl or wget is required" >&2
  exit 1
fi

chmod +x "$tmp"

if [ ! -d "$bin_dir" ] && [ -w "$(dirname "$bin_dir")" ]; then
  mkdir -p "$bin_dir"
fi

target="${bin_dir}/${binary}"
if [ -w "$bin_dir" ]; then
  mv "$tmp" "$target"
else
  if [ ! -d "$bin_dir" ]; then
    sudo mkdir -p "$bin_dir"
  fi
  sudo mv "$tmp" "$target"
fi

echo "installed ${target}"
echo "run: seekai --help"
