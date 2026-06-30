//go:build windows

package process

import (
	"os/exec"
	"syscall"
	"unsafe"
)

const createNoWindow = 0x08000000

var (
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	user32                = syscall.NewLazyDLL("user32.dll")
	getConsoleProcessList = kernel32.NewProc("GetConsoleProcessList")
	getConsoleWindow      = kernel32.NewProc("GetConsoleWindow")
	showWindow            = user32.NewProc("ShowWindow")
)

const swHide = 0

func HideWindow(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: createNoWindow,
	}
}

func HideConsoleWindowIfUnshared() {
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return
	}

	processes := make([]uint32, 8)
	count, _, _ := getConsoleProcessList.Call(uintptr(unsafe.Pointer(&processes[0])), uintptr(len(processes)))
	if count <= 1 {
		showWindow.Call(hwnd, swHide)
	}
}
