//go:build windows
// +build windows

package wallhaven

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
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

func KillDaemon(app *appcontext.AppContext) error {
	data, err := app.CacheTools.ReadLineFromFile("daemon.pid", 1)
	if err != nil {
		log.Fatalf("Failed to read PID file: %v", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return fmt.Errorf("invalid PID in file: %w", err)
	}

	cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("taskkill failed: %w", err)
	}

	app.CacheTools.DeleteFile("daemon.pid")

	return nil

}
