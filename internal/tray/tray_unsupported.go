//go:build !windows

package tray

import (
	"context"
	"fmt"

	"github.com/clooooud/wsl-tunneling/internal/config"
)

func Run(ctx context.Context, cfg config.Config, configPath string) error {
	return fmt.Errorf("tray is only supported on Windows")
}
