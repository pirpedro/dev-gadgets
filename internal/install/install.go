package install

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pirpedro/dev-gadgets/internal/catalog"
)

type Options struct {
	AssumeYes bool
}

func Install(ctx context.Context, it catalog.Item, opts Options) error {
	// Idempotency: verify first
	if it.Verify != "" {
		if ok := verify(it.Verify); ok {
			return nil
		}
	}

	// Escolhe estratégia baseada no SO e disponibilidade
	if it.Strategy.Release["url"] != "" {
		return runRelease(ctx, it)
	}

	// Detecta gerenciadores disponíveis
	has := func(bin string) bool {
		cmd := exec.Command("which", bin)
		return cmd.Run() == nil
	}

	// Instala uv em XDG se necessário
	installUvIfNeeded := func() error {
		if has("uv") {
			return nil
		}
		// Instalação simplificada: baixa binário para XDG
		xdgBin := getXdgBinDir()
		url := "https://github.com/astral-sh/uv/releases/latest/download/uv-x86_64-unknown-linux-musl.tar.gz"
		cmd := exec.Command("sh", "-c", fmt.Sprintf("curl -sSL %s | tar -xz -C %s", url, xdgBin))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("falha ao instalar uv: %v", err)
		}
		return nil
	}

	// Instala volta em XDG se necessário
	installVoltaIfNeeded := func() error {
		if has("volta") {
			return nil
		}
		xdgBin := getXdgBinDir()
		url := "https://get.volta.sh/latest"
		cmd := exec.Command("sh", "-c", fmt.Sprintf("curl %s | VOLTA_HOME=%s bash", url, xdgBin))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("falha ao instalar volta: %v", err)
		}
		return nil
	}

	// Python: uv/pipx
	if it.Strategy.Uv != "" {
		if opts.AssumeYes {
			// Modo automático: instala uv se não houver
			if has("uv") {
				return runUv(ctx, it.Strategy.Uv)
			} else {
				if err := installUvIfNeeded(); err != nil {
					return err
				}
				return runUv(ctx, it.Strategy.Uv)
			}
		} else {
			// Modo interativo: pergunta ao usuário
			var escolha string
			fmt.Printf("Você deseja instalar com uv para %s? (s/n): ", it.ID)
			fmt.Scanln(&escolha)
			if escolha == "s" {
				if !has("uv") {
					if err := installUvIfNeeded(); err != nil {
						return err
					}
				}
				return runUv(ctx, it.Strategy.Uv)
			}
		}
	}
	if it.Strategy.Pipx != "" {
		if opts.AssumeYes {
			if has("pipx") {
				return runPipx(ctx, it.Strategy.Pipx)
			}
		} else {
			var escolha string
			fmt.Printf("Você deseja instalar com pipx para %s? (s/n): ", it.ID)
			fmt.Scanln(&escolha)
			if escolha == "s" && has("pipx") {
				return runPipx(ctx, it.Strategy.Pipx)
			}
		}
	}

	// Node: volta/npm
	if it.Strategy.Volta != "" {
		if opts.AssumeYes {
			if has("volta") {
				return runVolta(ctx, it.Strategy.Volta)
			} else {
				if err := installVoltaIfNeeded(); err != nil {
					return err
				}
				return runVolta(ctx, it.Strategy.Volta)
			}
		} else {
			var escolha string
			fmt.Printf("Você deseja instalar com volta para %s? (s/n): ", it.ID)
			fmt.Scanln(&escolha)
			if escolha == "s" {
				if !has("volta") {
					if err := installVoltaIfNeeded(); err != nil {
						return err
					}
				}
				return runVolta(ctx, it.Strategy.Volta)
			}
		}
	}
	if it.Strategy.Npm != "" {
		if opts.AssumeYes {
			if has("npm") {
				return runNpm(ctx, it.Strategy.Npm)
			}
		} else {
			var escolha string
			fmt.Printf("Você deseja instalar com npm para %s? (s/n): ", it.ID)
			fmt.Scanln(&escolha)
			if escolha == "s" && has("npm") {
				return runNpm(ctx, it.Strategy.Npm)
			}
		}
	}

	// Demais gerenciadores
	if it.Strategy.Brew != "" && has("brew") {
		return runBrew(ctx, it.Strategy.Brew)
	}
	if it.Strategy.Apt != "" && has("apt-get") {
		return runApt(ctx, it.Strategy.Apt)
	}
	if it.Strategy.Dnf != "" && has("dnf") {
		return runDnf(ctx, it.Strategy.Dnf)
	}
	if it.Strategy.Pacman != "" && has("pacman") {
		return runPacman(ctx, it.Strategy.Pacman)
	}
	if it.Strategy.Zypper != "" && has("zypper") {
		return runZypper(ctx, it.Strategy.Zypper)
	}

	return fmt.Errorf("no viable strategy for %s", it.ID)
}

// Retorna o diretório XDG para binários
func getXdgBinDir() string {
	// Usa github.com/adrg/xdg
	// Se não conseguir, fallback para ~/.local/bin
	dir := ""
	if xdgBin, err := exec.Command("sh", "-c", "echo $XDG_DATA_HOME").Output(); err == nil && len(xdgBin) > 0 {
		dir = strings.TrimSpace(string(xdgBin)) + "/bin"
	} else {
		dir = fmt.Sprintf("%s/.local/bin", strings.TrimSpace(os.Getenv("HOME")))
	}
	// Verifica se está no PATH
	path := os.Getenv("PATH")
	if !strings.Contains(path, dir) {
		fmt.Fprintf(os.Stderr, "\n[dev-gadgets] Dica: adicione %s ao seu PATH para usar os binários instalados!\n", dir)
	}
	return dir
}

func verify(cmdline string) bool {
	parts := strings.Fields(cmdline)
	if len(parts) == 0 {
		return false
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	if err := cmd.Run(); err == nil {
		return true
	}
	return false
}
