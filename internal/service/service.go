package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/clooooud/wsl-tunneling/internal/config"
)

func Install(ctx context.Context, cfg config.Config, configPath string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("background installation is only supported on Windows")
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.LogDir, 0o755); err != nil {
		return err
	}
	logPath := filepath.Join(cfg.LogDir, "scheduled-task.log")
	command := fmt.Sprintf("powershell.exe -NoProfile -ExecutionPolicy Bypass -Command \"& '%s' daemon --config '%s' *> '%s'\"", exe, configPath, logPath)
	return runSchTasks(ctx, "/Create", "/F", "/SC", "ONLOGON", "/TN", cfg.ServiceName, "/TR", command)
}

func Uninstall(ctx context.Context, cfg config.Config) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("background installation is only supported on Windows")
	}
	return runSchTasks(ctx, "/Delete", "/F", "/TN", cfg.ServiceName)
}

func runSchTasks(ctx context.Context, args ...string) error {
	command := exec.CommandContext(ctx, "schtasks.exe", args...)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("schtasks.exe %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return nil
}
