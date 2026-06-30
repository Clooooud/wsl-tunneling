//go:build windows

package dns

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/clooooud/wsl-tunneling/internal/process"
)

func SearchSuffixes(ctx context.Context) ([]string, error) {
	script := strings.Join([]string{
		"$items = @()",
		"$global = (Get-DnsClientGlobalSetting).SuffixSearchList",
		"if ($global) { $items += $global }",
		"Get-DnsClient | ForEach-Object { if ($_.ConnectionSpecificSuffix) { $items += $_.ConnectionSpecificSuffix } }",
		"$items | Where-Object { $_ } | Select-Object -Unique",
	}, "; ")
	command := exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script)
	process.HideWindow(command)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return nil, err
	}
	return NormalizeSearchSuffixes(strings.Split(stdout.String(), "\n")), nil
}
