package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/clooooud/wsl-tunneling/internal/config"
	"github.com/clooooud/wsl-tunneling/internal/process"
)

func Install(ctx context.Context, cfg config.Config, configPath string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("background installation is only supported on Windows")
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	return runReg(ctx, "add", autostartKey, "/v", autostartName, "/t", "REG_SZ", "/d", autostartCommand(exe, configPath), "/f")
}

func IsInstalled(ctx context.Context, cfg config.Config) (bool, error) {
	if runtime.GOOS != "windows" {
		return false, fmt.Errorf("background installation is only supported on Windows")
	}
	command := exec.CommandContext(ctx, "reg.exe", "query", autostartKey, "/v", autostartName)
	process.HideWindow(command)
	output, err := command.CombinedOutput()
	if err == nil {
		return true, nil
	}
	text := strings.ToLower(string(output))
	if strings.Contains(text, "unable to find") || strings.Contains(text, "the system was unable to find") || strings.Contains(text, "introuvable") {
		return false, nil
	}
	return false, fmt.Errorf("reg.exe query %s /v %s: %w: %s", autostartKey, autostartName, err, strings.TrimSpace(string(output)))
}

func Uninstall(ctx context.Context, cfg config.Config) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("background installation is only supported on Windows")
	}
	installed, err := IsInstalled(ctx, cfg)
	if err != nil || !installed {
		return err
	}
	return runReg(ctx, "delete", autostartKey, "/v", autostartName, "/f")
}

const autostartKey = `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
const autostartName = config.DefaultServiceName

func autostartCommand(exe string, configPath string) string {
	command := fmt.Sprintf("\"%s\"", exe)
	if configPath != "" && configPath != config.DefaultPath() {
		command = fmt.Sprintf("%s --config \"%s\"", command, configPath)
	}
	return command
}

func runReg(ctx context.Context, args ...string) error {
	command := exec.CommandContext(ctx, "reg.exe", args...)
	process.HideWindow(command)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("reg.exe %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return nil
}
