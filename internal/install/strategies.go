package install

import (
	"context"
	"fmt"
	"os/exec"
)

func runUv(ctx context.Context, pkg string) error {
	cmd := exec.CommandContext(ctx, "uv", "pip", "install", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("uv install failed: %v\n%s", err, out)
	}
	return nil
}

func runNpm(ctx context.Context, pkg string) error {
	cmd := exec.CommandContext(ctx, "npm", "install", "-g", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm install failed: %v\n%s", err, out)
	}
	return nil
}

func runVolta(ctx context.Context, pkg string) error {
	cmd := exec.CommandContext(ctx, "volta", "install", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("volta install failed: %v\n%s", err, out)
	}
	return nil
}

func runBrew(ctx context.Context, pkg string) error {
	cmd := exec.CommandContext(ctx, "brew", "install", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("brew install failed: %v\n%s", err, out)
	}
	return nil
}

func runApt(ctx context.Context, pkg string) error {
	// TODO: sudo apt-get update && sudo apt-get install -y pkg
	cmd := exec.CommandContext(ctx, "sudo", "apt-get", "update")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("apt-get update failed: %v\n%s", err, out)
	}
	cmd = exec.CommandContext(ctx, "sudo", "apt-get", "install", "-y", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("apt-get install failed: %v\n%s", err, out)
	}
	return nil
}

func runPipx(ctx context.Context, pkg string) error {
	// TODO: pipx install pkg
	cmd := exec.CommandContext(ctx, "pipx", "install", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pipx install failed: %v\n%s", err, out)
	}
	return nil
}

func runDnf(ctx context.Context, pkg string) error {
	cmd := exec.CommandContext(ctx, "sudo", "dnf", "install", "-y", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dnf install failed: %v\n%s", err, out)
	}
	return nil
}

func runPacman(ctx context.Context, pkg string) error {
	cmd := exec.CommandContext(ctx, "sudo", "pacman", "-Sy", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pacman install failed: %v\n%s", err, out)
	}
	return nil
}

func runZypper(ctx context.Context, pkg string) error {
	cmd := exec.CommandContext(ctx, "sudo", "zypper", "install", "-y", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("zypper install failed: %v\n%s", err, out)
	}
	return nil
}
