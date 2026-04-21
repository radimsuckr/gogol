//go:build !windows

package main

// enableANSI is a no-op on Unix-like systems where ANSI is supported by default
func enableANSI() error {
	return nil
}
