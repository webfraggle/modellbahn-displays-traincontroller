//go:build !windows

package console

// Attach is a no-op on macOS/Linux, where stdout is always connected.
func Attach() {}
