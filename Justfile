# =============================================================================
# 📚 Documentation
# =============================================================================
# This justfile provides a comprehensive build system for Go projects of any size.
# It supports development, testing, building, and deployment workflows.
#
# Quick Start:
# 1. Install 'just': https://github.com/casey/just
# 2. Copy this justfile to your project root
# 3. Run `just init` to initialize the project
# 4. Run `just --list` to see available commands
#
# Configuration:
# The justfile can be configured in several ways (in order of precedence):
# 1. Command line: just GOOS=darwin build
# 2. Environment variables: export GOOS=darwin
# 3. .env file in project root
# 4. Default values in this justfile

# =============================================================================
# Variables
# =============================================================================

# You may override via env/.env if needed
go          := `command -v go || printf :`

export GOPATH := env_var_or_default("GOPATH", `command -v go >/dev/null 2>&1 && go env GOPATH || echo ""`)
export GOOS   := env_var_or_default("GOOS",   `command -v go >/dev/null 2>&1 && go env GOOS || echo ""`)
export GOARCH := env_var_or_default("GOARCH", `command -v go >/dev/null 2>&1 && go env GOARCH || echo ""`)
export CGO_ENABLED := env_var_or_default("CGO_ENABLED", "0")

gobin       := GOPATH + "/bin"

# ldflags embedding version info from the base (version/git_commit/git_branch/build_time/build_by)
# Uses current module path when available; falls back gracefully.
module      := `go list -m 2>/dev/null || echo ""`
ld_flags    := if module != "" {
  "-s -w \
   -X '" + module + "/pkg/version.Version=" + version + "' \
   -X '" + module + "/pkg/version.Commit=" + git_commit + "' \
   -X '" + module + "/pkg/version.Branch=" + git_branch + "' \
   -X '" + module + "/pkg/version.BuildTime=" + build_time + "' \
   -X '" + module + "/pkg/version.BuildBy=" + build_by + "'"
} else {
  "-s -w"
}

# Optional build tags (override via env if needed)
build_tags  := env_var_or_default("GO_BUILD_TAGS", "")
all_tags    := build_tags

# Test settings
test_timeout := env_var_or_default("GO_TEST_TIMEOUT", "5m")
bench_time   := env_var_or_default("GO_BENCH_TIME", "2s")

# =============================================================================
# Default Recipe
# =============================================================================

# Show available recipes in fzf format
default:
  @just --choose

# =============================================================================
# Project Setup
# =============================================================================
# Install deps
[group('🛠️  Development')]
@install:
  {{go}} mod download
  # Optional: install dev tools if present in go.work/tools or similar (no-op if missing)
  # {{go}} install github.com/golangci/golangci-lint/cmd/golangci-lint@latest || true
  # {{go}} install mvdan.cc/gofumpt@latest || true
  # {{go}} install golang.org/x/vuln/cmd/govulncheck@latest || true
  # {{go}} install github.com/golang/mock/mockgen@latest || true
  # {{go}} install github.com/air-verse/air@latest || true

# test
[group('🧪 Quality')]
test *ARGS:
  {{go}} test -v -race -timeout {{test_timeout}} -cover ./... {{ARGS}}

# lint
[group('🧪 Quality')]
lint:
  if command -v golangci-lint >/dev/null 2>&1; then \
    golangci-lint run --fix || true; \
  else \
    {{go}} vet ./...; \
  fi

# fmt
[group('🧪 Quality')]
fmt:
  {{go}} fmt ./...
  if command -v gofumpt >/dev/null 2>&1; then gofumpt -l -w .; fi

# build
[group('🛠️  Development')]
build:
  mkdir -p "{{bin_dir}}"
  {{go}} build -tags '{{build_tags}}' -ldflags '{{ld_flags}}' -o "{{bin_dir}}/{{project_name}}" {{entrypoint}}

[group('🛠️  Development')]
generate:
  {{go}} generate ./...

[group('🛠️  Development')]
test-coverage:
  {{go}} test -v -race -coverprofile=coverage.out ./...
  {{go}} tool cover -html=coverage.out -o coverage.html
  @echo "→ open coverage.html"

[group('🧪 Quality')]
bench:
  {{go}} test -bench=. -benchmem -run=^$ -benchtime={{bench_time}} ./...

[group('🧪 Quality')]
vet:
  {{go}} vet ./...

[group('🧪 Quality')]
security:
  if command -v govulncheck >/dev/null 2>&1; then govulncheck ./...; else echo "install govulncheck first"; fi

[group('🛠️  Development')]
build-all:
  #!/usr/bin/env bash
  set -e
  mkdir -p "{{dist_dir}}"
  platforms=("linux/amd64/-" "linux/arm64/-" "linux/arm/7" "darwin/amd64/-" "darwin/arm64/-" "windows/amd64/-" "windows/arm64/-")
  for p in "${platforms[@]}"; do
    os="${p%%/*}"; rest="${p#*/}"; arch="${rest%%/*}"; arm="${rest##*/}"
    ext=""; [ "$os" = "windows" ] && ext=".exe"
    out="{{dist_dir}}/{{project_name}}-${os}-${arch}${ext}"
    if [ "$arm" != "-" ]; then export GOARM="$arm"; fi
    GOOS="$os" GOARCH="$arch" CGO_ENABLED="{{CGO_ENABLED}}" {{go}} build -tags '{{build_tags}}' -ldflags '{{ld_flags}}' -o "$out" {{entrypoint}}
    (cd "{{dist_dir}}" && tar -czf "$(basename "$out").tar.gz" "$(basename "$out")")
    rm -f "$out"
  done
  echo "Artifacts ready at {{dist_dir}}"

[group('🛠️  Development')]
docs:
  mkdir -p "{{docs_dir}}"
  {{go}} doc -all > "{{docs_dir}}/API.md"

[group('🛠️  Development')]
init:
  @echo "Initializing project..."
  @mkdir -p {{bin_dir}} {{dist_dir}} {{docs_dir}}
  @$(go) mod tidy
  @$(go) mod vendor
  @echo "Project initialized."


gadgets *ARGS:
  @just dev {{ARGS}}

# Run dev-gadgets in interactive install mode
[group('🛠️  Development')]
gadgets-interactive:
  @just gadgets install --interactive

# Run dev-gadgets with all available commands
[group('🛠️  Development')]
gadgets-install-all:
  @just gadgets install --all

import? "~/.config/just/Justfile"
