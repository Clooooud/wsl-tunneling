package daemon

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/clooooud/wsl-tunneling/internal/config"
	"github.com/clooooud/wsl-tunneling/internal/network"
)

func Run(ctx context.Context, cfg config.Config, output io.Writer) error {
	manager := network.NewManager(cfg)
	if err := manager.Start(ctx); err != nil {
		return err
	}
	fmt.Fprintf(output, "wsl-tunneling daemon started for distro %s\n", cfg.Distro)

	intervalSeconds := cfg.SupervisorInterval
	if intervalSeconds <= 0 {
		intervalSeconds = 10
	}
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintln(output, "stopping wsl-tunneling daemon")
			return manager.Stop(context.Background())
		case <-ticker.C:
			status, err := manager.Status(ctx)
			if err != nil || !status.ForwarderUp {
				fmt.Fprintf(output, "forwarder unhealthy, restarting: %v\n", err)
				_ = manager.Stop(ctx)
				if startErr := manager.Start(ctx); startErr != nil {
					fmt.Fprintf(output, "restart failed: %v\n", startErr)
				}
			}
		}
	}
}
