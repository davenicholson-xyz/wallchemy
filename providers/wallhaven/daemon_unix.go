//go:build !windows
// +build !windows

package wallhaven

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
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
	defer nullFile.Close()

	cmd.Stdin = nullFile
	cmd.Stdout = nullFile
	cmd.Stderr = nullFile

	// Properly detach the process (Unix)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	// Start the daemon process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	// Optionally log the PID of the daemon
	fmt.Printf("Daemon started with PID %d\n", cmd.Process.Pid)

	// ✅ Immediately exit the parent process cleanly
	os.Exit(0)

	return nil // unreachable but included for completeness

	// nullFile, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	// if err != nil {
	// 	return fmt.Errorf("failed to open null device: %w", err)
	// }
	// cmd.Stdin = nullFile
	// cmd.Stdout = nullFile
	// cmd.Stderr = nullFile
	//
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Setsid: true,
	// }
	//
	// return cmd.Start()
}
