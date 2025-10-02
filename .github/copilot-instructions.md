# Copilot Instructions — git-gadgets (Go)

Objetivo: migrar de scripts shell para um binário único em Go (Linux/macOS) com TUI Bubble Tea, modo não-interativo para CI e instaladores idempotentes dirigidos por catálogo.

## TL;DR (produto e UX)

- Comandos: `gadgets install [--all|--interactive] [--only item1,item2] [--yes] [--dry-run]`, `list`, `doctor`, `update`.
- TUI (Bubble Tea): multi-seleção com ↑/↓ mover, espaço marcar, enter confirmar; mostra plano, confirma (exceto com `--yes`), executa com progresso e resume sucesso/pulo/falha.
- Não-interativo: `--all` (curado), `--only`, `--dry-run`, `--yes` (sem prompts).

## Arquitetura (Go 1.22+)

- CLI: `spf13/cobra` em `cmd/dev-gadgets` (subcomandos em `cmd/dev-gadgets/cmd/…`).
- TUI: `charmbracelet/bubbletea`, `bubbles` (lista/multi), `lipgloss` (estilo). UI principal em `internal/ui/select_items.go` com `SelectItemsModel` (campos: `items []catalog.Item`, `selected map[string]bool`, keymap ↑/↓/space/enter/quit, `Update`/`View`).
- Catálogo: YAML em `config/catalog.yaml` (schema: id, name, verify, strategies[brew|apt|dnf|pacman|zypper|pipx|release]). Código em `internal/catalog`.
- Instaladores: estratégias ordenadas por item; escolhe a viável, instala e valida com `--version`; idempotente (se ok, pula). Execução paralela com `errgroup` onde seguro.
- Paths/estado: `github.com/adrg/xdg` p/ `~/.local/share/dev-gadgets` e instalar binários em `~/.local/bin` (avisar se não estiver no PATH).

## Fluxos de dev

- Build/test: `go build ./...`, `go test ./...`. Lint opcional (`go vet`, golangci-lint se adicionado).
- Execução rápida: `go run ./cmd/dev-gadgets --help`; para TUI: `go run ./cmd/dev-gadgets install --interactive`.
- Releases: GoReleaser gera assets ZIP para linux/darwin (amd64/arm64) + script instalador 1-liner que baixa o binário correto para `~/.local/bin` (reutilizar/atualizar `install.sh`).
- Devcontainer/compose: úteis para testar ambientes de pacote (apt/dnf). Ajuste paths se necessário.

## Convenções

- Padrão Cobra: um pacote por comando em `cmd/dev-gadgets/cmd/...` com `RunE` e flags (`--all`, `--only`, `--yes`, `--dry-run`, `--interactive`).
- UI isolada em `internal/ui`; lógica de instalação em `internal/install` (ex.: `strategies.go`, `brew.go`, `apt.go`, `release.go`).
- Catálogo declarativo: exemplo mínimo em `config/catalog.yaml`:
  - `- id: git-town; name: Git Town; verify: git-town --version; strategies: [brew: git-town, apt: git-town, release: {url: ..., bin: git-town}]`.
- Logs: stdlib (`log`/`slog`), mensagens amigáveis; erros agregados por item.

## Migração do estado atual

- Hoje: `bin/`, `helper/`, `Makefile` montam scripts shell; manter `install.sh` (ajustado) para baixar a release Go. Portar funcionalidades para Cobra/Bubble Tea seguindo estruturas acima. Compleções passam a ser geradas via Cobra.

## Gotchas

- Idempotência: sempre verificar presença/versão antes de instalar; em CI, combinar `--yes` + `--dry-run` para planejar.
- Estratégias com sudo (apt/dnf/pacman/zypper): detectar pacote antes e informar necessidade de privilégio.
- PATH: garantir `~/.local/bin` no PATH; imprimir dica se ausente.
