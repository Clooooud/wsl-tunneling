package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultsValidateAfterDistro(t *testing.T) {
	cfg := Defaults()
	cfg.Distro = "Ubuntu"
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidateRequiresDistro(t *testing.T) {
	cfg := Defaults()
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want distro error")
	}
}

func TestSaveWritesConfigJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.json")
	cfg := Defaults()
	cfg.Distro = "Ubuntu-24.04"

	if err := Save(path, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.HasSuffix(string(content), "\n") {
		t.Fatalf("saved config should end with newline: %q", content)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.Distro != cfg.Distro {
		t.Fatalf("Load().Distro = %q, want %q", loaded.Distro, cfg.Distro)
	}
}
