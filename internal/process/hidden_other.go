//go:build !windows

package process

import "os/exec"

func HideWindow(command *exec.Cmd) {}

func HideConsoleWindowIfUnshared() {}
