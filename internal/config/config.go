package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	DefaultGvisorVersion = "v0.8.9"
	DefaultInterface     = "wsl-tunneling"
	DefaultSubnet        = "192.168.127.0/24"
	DefaultGatewayIP     = "192.168.127.1"
	DefaultDeviceIP      = "192.168.127.2"
	DefaultStateDirWSL   = "/mnt/wsl/wsl-tunneling"
	DefaultServiceName   = "wsl-tunneling"
)

type Config struct {
	Distro              string      `json:"distro"`
	GvisorVersion       string      `json:"gvisorVersion"`
	CacheDir            string      `json:"cacheDir"`
	StateDir            string      `json:"stateDir"`
	LogDir              string      `json:"logDir"`
	StateDirWSL         string      `json:"stateDirWsl"`
	InterfaceName       string      `json:"interfaceName"`
	Subnet              string      `json:"subnet"`
	GatewayIP           string      `json:"gatewayIp"`
	DeviceIP            string      `json:"deviceIp"`
	DNSSearchSuffixes   []string    `json:"dnsSearchSuffixes,omitempty"`
	GVProxyURL          string      `json:"gvproxyUrl,omitempty"`
	GVForwarderURL      string      `json:"gvforwarderUrl,omitempty"`
	GitHubAPIBaseURL    string      `json:"githubApiBaseUrl"`
	ServiceName         string      `json:"serviceName"`
	SupervisorInterval  int         `json:"supervisorIntervalSeconds"`
	TerminateOnStop     bool        `json:"terminateOnStop"`
	DisableAutoResolv   bool        `json:"disableAutoResolv"`
	DownloadPermissions os.FileMode `json:"-"`
}

func Defaults() Config {
	dataRoot := defaultDataRoot()

	return Config{
		GvisorVersion:       DefaultGvisorVersion,
		CacheDir:            filepath.Join(dataRoot, "bin"),
		StateDir:            filepath.Join(dataRoot, "state"),
		LogDir:              filepath.Join(dataRoot, "logs"),
		StateDirWSL:         DefaultStateDirWSL,
		InterfaceName:       DefaultInterface,
		Subnet:              DefaultSubnet,
		GatewayIP:           DefaultGatewayIP,
		DeviceIP:            DefaultDeviceIP,
		GitHubAPIBaseURL:    "https://api.github.com",
		ServiceName:         DefaultServiceName,
		SupervisorInterval:  10,
		TerminateOnStop:     false,
		DisableAutoResolv:   true,
		DownloadPermissions: 0o755,
	}
}

func DefaultPath() string {
	return filepath.Join(defaultConfigRoot(), "config.json")
}

func Load(path string) (Config, error) {
	cfg := Defaults()
	if path == "" {
		path = DefaultPath()
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config: %w", err)
	}

	if err := json.Unmarshal(content, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}

func SaveExample(path string) error {
	cfg := Defaults()
	cfg.Distro = "Ubuntu"
	return Save(path, cfg)
}

func Save(path string, cfg Config) error {
	if path == "" {
		path = DefaultPath()
	}
	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.Distro) == "" {
		return errors.New("distro is required; pass --distro or set distro in config.json")
	}
	if strings.TrimSpace(cfg.GvisorVersion) == "" {
		return errors.New("gvisorVersion is required")
	}
	if _, err := netip.ParsePrefix(cfg.Subnet); err != nil {
		return fmt.Errorf("invalid subnet %q: %w", cfg.Subnet, err)
	}
	if _, err := netip.ParseAddr(cfg.GatewayIP); err != nil {
		return fmt.Errorf("invalid gatewayIp %q: %w", cfg.GatewayIP, err)
	}
	if _, err := netip.ParseAddr(cfg.DeviceIP); err != nil {
		return fmt.Errorf("invalid deviceIp %q: %w", cfg.DeviceIP, err)
	}
	if strings.TrimSpace(cfg.InterfaceName) == "" {
		return errors.New("interfaceName is required")
	}
	if strings.TrimSpace(cfg.StateDirWSL) == "" || !strings.HasPrefix(cfg.StateDirWSL, "/") {
		return errors.New("stateDirWsl must be an absolute Linux path")
	}
	if cfg.SupervisorInterval < 1 {
		return errors.New("supervisorIntervalSeconds must be greater than zero")
	}
	return nil
}

func EnsureDirs(cfg Config) error {
	for _, path := range []string{cfg.CacheDir, cfg.StateDir, cfg.LogDir} {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return fmt.Errorf("create %s: %w", path, err)
		}
	}
	return nil
}

func defaultConfigRoot() string {
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "wsl-tunneling")
		}
	}
	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		return filepath.Join(configHome, "wsl-tunneling")
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config", "wsl-tunneling")
	}
	return "."
}

func defaultDataRoot() string {
	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return filepath.Join(localAppData, "wsl-tunneling")
		}
	}
	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		return filepath.Join(dataHome, "wsl-tunneling")
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".local", "share", "wsl-tunneling")
	}
	return "."
}
