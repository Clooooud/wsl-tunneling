package wsl

import "testing"

func TestWindowsPathToWSL(t *testing.T) {
	got, err := WindowsPathToWSL(`C:\Users\me\bin\gvproxy.exe`)
	if err != nil {
		t.Fatalf("WindowsPathToWSL() error = %v", err)
	}
	want := "/mnt/c/Users/me/bin/gvproxy.exe"
	if got != want {
		t.Fatalf("WindowsPathToWSL() = %q, want %q", got, want)
	}
}

func TestWindowsPathToWSLRejectsRelative(t *testing.T) {
	if _, err := WindowsPathToWSL(`bin\gvproxy.exe`); err == nil {
		t.Fatal("WindowsPathToWSL() error = nil, want error")
	}
}
