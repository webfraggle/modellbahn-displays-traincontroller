//go:build windows

package console

import (
	"os"
	"syscall"
)

const attachParentProcess = uintptr(0xFFFFFFFF) // ATTACH_PARENT_PROCESS ((DWORD)-1)

// Attach connects stdout/stderr to the console of the calling process.
// The Windows binary uses the GUI subsystem and has no console of its own,
// so without this nothing printed would be visible in cmd/PowerShell.
func Attach() {
	if h := os.Stdout; h != nil && h.Fd() != 0 && syscall.Handle(h.Fd()) != syscall.InvalidHandle {
		return // stdout already usable (e.g. redirected to a file or pipe)
	}
	attach := syscall.NewLazyDLL("kernel32.dll").NewProc("AttachConsole")
	if r, _, _ := attach.Call(attachParentProcess); r == 0 {
		return
	}
	if f, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0); err == nil {
		os.Stdout = f
		os.Stderr = f
	}
}
