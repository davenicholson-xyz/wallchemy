//go:build !windows
// +build !windows

package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/davenicholson-xyz/wallchemy/app"
)

func LaunchDaemon(appCtx *app.AppContext) (string, error) {

	execPath, err := os.Executable()
	if err != nil {
		app.Logger.WithField("execPath", execPath).Debug("Found executable path")
		return "", fmt.Errorf("could not determine executable path: %w", err)
	}

	cmd := exec.Command(execPath, "--launchdaemon", "--port", strconv.Itoa(appCtx.Config.Port))
	nullFile, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return "", fmt.Errorf("failed to open null device: %w", err)
	}
	defer nullFile.Close()

	cmd.Stdin = nullFile
	cmd.Stdout = nullFile
	cmd.Stderr = nullFile

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start daemon: %w", err)
	}

	os.Exit(0)

	return "", nil

}
