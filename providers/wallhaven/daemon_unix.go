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
	"github.com/davenicholson-xyz/wallchemy/logger"
)

func LaunchDaemon() error {

	execPath, err := os.Executable()
	if err != nil {
		logger.Log.WithField("execPath", execPath).Debug("Found executable path")
		return fmt.Errorf("could not determine executable path: %w", err)
	}

	cmd := exec.Command(execPath, "-startdaemon")
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
		logger.Log.Debug("No daemon PID found")
		return nil
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		logger.Log.WithField("PID", pid).Debug("Invalid PID found")
		return nil
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		logger.Log.WithField("PID", pid).Debug("Unable to find process with PID")
		return nil
	}

	if err := process.Kill(); err != nil {
		logger.Log.WithError(err).Debug("Could not kill process")
		return nil
	}

	app.CacheTools.DeleteFile("daemon.pid")

	logger.Log.Info("Daemon process stopped")

	return nil

}
