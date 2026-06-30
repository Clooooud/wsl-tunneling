package main

import (
	"os"

	"github.com/clooooud/wsl-tunneling/internal/cli"
	"github.com/clooooud/wsl-tunneling/internal/process"
)

func main() {
	if cli.StartsTray(os.Args[1:]) {
		process.HideConsoleWindowIfUnshared()
	}
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
