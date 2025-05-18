//go:build !windows
// +build !windows

package wallhaven

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
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

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	fmt.Printf("Daemon started with PID %d\n", cmd.Process.Pid)

	os.Exit(0)

	return nil

}

func KillDaemon(app *appcontext.AppContext) error {
	data, err := app.CacheTools.ReadLineFromFile("daemon.pid", 1)
	if err != nil {
		return nil
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return nil
		// return fmt.Errorf("invalid PID in file: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return nil
		// return fmt.Errorf("failed to find process with PID %d: %w", pid, err)
	}

	if err := process.Kill(); err != nil {
		return nil
		// return fmt.Errorf("failed to kill process: %w", err)
	}

	app.CacheTools.DeleteFile("daemon.pid")

	return nil

}
