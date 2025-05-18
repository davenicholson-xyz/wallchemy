//go:build windows
// +build windows

package wallhaven

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const (
	CREATE_NEW_PROCESS_GROUP = 0x00000200
	CREATE_NO_WINDOW         = 0x08000000
)

func LaunchDaemon() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine executable path: %w", err)
	}

	cmd := exec.Command(execPath)
	cmd.Env = append(os.Environ(), "WALLCHEMY_STARTDAEMON=1")

	nullFile, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open null device: %w", err)
	}
	cmd.Stdin = nullFile
	cmd.Stdout = nullFile
	cmd.Stderr = nullFile

	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: CREATE_NEW_PROCESS_GROUP | CREATE_NO_WINDOW,
	}

	return cmd.Start()
}
