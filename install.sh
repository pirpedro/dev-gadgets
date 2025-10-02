#!/usr/bin/env bash
set -euo pipefail

# One-liner installer for dev-gadgets (Linux/macOS, amd64/arm64)
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/pirpedro/dev-gadgets/main/install.sh | bash
# Options:
#   PREFIX=~/.local/bin   # install destination (default: ~/.local/bin)
#   VERSION=vX.Y.Z        # specific tag (default: latest)

REPO="pirpedro/dev-gadgets"
BIN="dev-gadgets"

log() { printf "[dev-gadgets] %s\n" "$*"; }
err() { printf "[dev-gadgets] ERROR: %s\n" "$*" >&2; exit 1; }

detect_os_arch() {
  local os arch
  case "$(uname -s)" in
    Linux) os=linux ;;
    Darwin) os=darwin ;;
    *) err "unsupported OS: $(uname -s)" ;;
  esac
  case "$(uname -m)" in
    x86_64|amd64) arch=amd64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) err "unsupported ARCH: $(uname -m)" ;;
  esac
  echo "$os" "$arch"
}

ensure_cmd() { command -v "$1" >/dev/null 2>&1 || err "missing dependency: $1"; }

latest_tag() {
  # Prefer GitHub API. Fallback to release redirect if curl without API tokens.
  local tag
  tag=$(curl -fsSL https://api.github.com/repos/$REPO/releases/latest | sed -n 's/.*"tag_name" *: *"\([^"]*\)".*/\1/p' | head -n1 || true)
  if [ -z "$tag" ]; then
    # Follow redirect to /releases/tag/<tag>
    tag=$(curl -fsIL -o /dev/null -w '%{url_effective}' https://github.com/$REPO/releases/latest | sed 's#.*/tag/##')
  fi
  [ -n "$tag" ] || err "could not resolve latest tag"
  echo "$tag"
}

install_dir() {
  local p="${PREFIX:-$HOME/.local/bin}"
  mkdir -p "$p" || true
  echo "$p"
}

in_path() {
  case ":$PATH:" in *":$1:"*) return 0;; esac; return 1
}

main() {
  ensure_cmd curl
  local os arch; read -r os arch < <(detect_os_arch)
  local tag="${VERSION:-}"
  if [ -z "$tag" ]; then tag=$(latest_tag); fi
  local name="${BIN}_${os}_${arch}"
  local url_zip="https://github.com/${REPO}/releases/download/${tag}/${name}.zip"
  local url_tgz="https://github.com/${REPO}/releases/download/${tag}/${name}.tar.gz"

  local tmp
  tmp=$(mktemp -d -t dev-gadgets.XXXXXX)
  trap 'rm -rf "$tmp"' EXIT

  log "installing ${BIN} ${tag} for ${os}/${arch}"

  # Try ZIP first, then TGZ
  local archive="${tmp}/pkg"
  if curl -fsSL "$url_zip" -o "$archive.zip"; then
    if command -v unzip >/dev/null 2>&1; then
      unzip -qo "$archive.zip" -d "$tmp"
    elif command -v bsdtar >/dev/null 2>&1; then
      bsdtar -xf "$archive.zip" -C "$tmp"
    else
      err "need 'unzip' or 'bsdtar' to extract zip"
    fi
  elif curl -fsSL "$url_tgz" -o "$archive.tgz"; then
    tar -xzf "$archive.tgz" -C "$tmp"
  else
    err "could not download release asset: $url_zip or $url_tgz"
  fi

  # Find binary in extracted files
  local binpath
  binpath=$(find "$tmp" -type f -name "$BIN" -perm -u+x -print -quit)
  if [ -z "$binpath" ]; then
    # maybe binary not marked executable inside archive; try any matching name
    binpath=$(find "$tmp" -type f -name "$BIN" -print -quit)
  fi
  [ -n "$binpath" ] || err "binary '$BIN' not found in archive"

  local dest; dest=$(install_dir)
  if cp "$binpath" "$dest/$BIN" 2>/dev/null; then
    :
  else
    log "need sudo to write to $dest"
    ensure_cmd sudo
    sudo cp "$binpath" "$dest/$BIN"
  fi
  chmod +x "$dest/$BIN" || true

  if ! in_path "$dest"; then
    printf "\n[dev-gadgets] Hint: add %s to your PATH (e.g., add to shell rc)\n\n" "$dest"
  fi

  "$dest/$BIN" --version || true
  log "installed at $dest/$BIN"
}

main "$@"
