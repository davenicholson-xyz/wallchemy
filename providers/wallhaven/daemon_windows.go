//go:build windows
// +build windows

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

const (
	CREATE_NEW_PROCESS_GROUP = 0x00000200
	CREATE_NO_WINDOW         = 0x08000000
)

func LaunchDaemon() error {
	execPath, err := os.Executable()
	if err != nil {
		logger.Log.WithField("execPath", execPath).Debug("Found executable path")
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

	cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
	if err := cmd.Run(); err != nil {
		logger.Log.WithError(err).Debug("Could not kill process")
		return nil
	}

	app.CacheTools.DeleteFile("daemon.pid")

	logger.Log.Info("Daemon process stopped")

	return nil

}
