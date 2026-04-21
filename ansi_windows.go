//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

// enableANSI enables ANSI escape code support on Windows
func enableANSI() error {
	// Windows console API constants
	const (
		STD_OUTPUT_HANDLE                  = ^uintptr(10) // -11
		ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	)

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getStdHandle := kernel32.NewProc("GetStdHandle")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	// Get stdout handle
	handle, _, _ := getStdHandle.Call(STD_OUTPUT_HANDLE)
	if handle == 0 {
		return fmt.Errorf("failed to get stdout handle")
	}

	// Get current console mode
	var mode uint32
	ret, _, err := getConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		return fmt.Errorf("GetConsoleMode failed: %v", err)
	}

	// Enable virtual terminal processing
	mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	ret, _, err = setConsoleMode.Call(handle, uintptr(mode))
	if ret == 0 {
		return fmt.Errorf("SetConsoleMode failed: %v", err)
	}

	return nil
}
