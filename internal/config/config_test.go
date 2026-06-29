package config

import "testing"

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
