//go:build windows
// +build windows

package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/davenicholson-xyz/wallchemy/app"
)

const (
	CREATE_NEW_PROCESS_GROUP = 0x00000200
	CREATE_NO_WINDOW         = 0x08000000
)

func LaunchDaemon(appCtx *app.AppContext) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not determine executable path: %w", err)
	}
	app.Logger.WithField("execPath", execPath).Debug("Found executable path")

	cmd := exec.Command(execPath, "--launchdaemon", "--port", strconv.Itoa(appCtx.Config.Port))
	nullFile, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return "", fmt.Errorf("failed to open null device: %w", err)
	}
	cmd.Stdin = nullFile
	cmd.Stdout = nullFile
	cmd.Stderr = nullFile

	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: CREATE_NEW_PROCESS_GROUP | CREATE_NO_WINDOW,
	}

	return "", cmd.Start()
}
