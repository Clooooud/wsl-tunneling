package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTrayExecutablePrefersLocalGuiSibling(t *testing.T) {
	dir := t.TempDir()
	tray := filepath.Join(dir, "wsl-tunneling.exe")
	if err := os.WriteFile(tray, []byte{}, 0o755); err != nil {
		t.Fatal(err)
	}

	cli := filepath.Join(dir, "wsl-tunneling-cli.exe")
	if got := trayExecutable(cli); got != tray {
		t.Fatalf("trayExecutable(%q) = %q, want %q", cli, got, tray)
	}
}

func TestTrayExecutablePrefersReleaseGuiSibling(t *testing.T) {
	dir := t.TempDir()
	tray := filepath.Join(dir, "wsl-tunneling-windows-amd64.exe")
	if err := os.WriteFile(tray, []byte{}, 0o755); err != nil {
		t.Fatal(err)
	}

	cli := filepath.Join(dir, "wsl-tunneling-cli-windows-amd64.exe")
	if got := trayExecutable(cli); got != tray {
		t.Fatalf("trayExecutable(%q) = %q, want %q", cli, got, tray)
	}
}

func TestTrayExecutableKeepsOriginalWithoutGuiSibling(t *testing.T) {
	dir := t.TempDir()
	cli := filepath.Join(dir, "wsl-tunneling-cli.exe")
	if got := trayExecutable(cli); got != cli {
		t.Fatalf("trayExecutable(%q) = %q, want %q", cli, got, cli)
	}
}
