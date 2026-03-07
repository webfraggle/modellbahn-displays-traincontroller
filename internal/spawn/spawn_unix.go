//go:build !windows

package spawn

import (
	"os"
	"os/exec"
	"syscall"
)

// Detach spawns a copy of the current process with --bg added to the args,
// in a new session (detached from the calling process), then exits immediately.
func Detach() {
	args := append([]string{"--bg"}, os.Args[1:]...)
	cmd := exec.Command(os.Args[0], args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	_ = cmd.Start()
	os.Exit(0)
}
