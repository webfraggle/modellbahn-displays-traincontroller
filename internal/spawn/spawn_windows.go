//go:build windows

package spawn

import (
	"os"
	"os/exec"
	"syscall"
)

// Detach spawns a copy of the current process with --bg added to the args,
// fully detached from the caller (no console window), then exits immediately.
func Detach() {
	args := append([]string{"--bg"}, os.Args[1:]...)
	cmd := exec.Command(os.Args[0], args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
		HideWindow:    true,
	}
	_ = cmd.Start()
	os.Exit(0)
}
