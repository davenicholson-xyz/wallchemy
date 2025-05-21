//go:build !windows
// +build !windows

package wallhaven

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/logger"
)

func LaunchDaemon(app *appcontext.AppContext) error {

	execPath, err := os.Executable()
	if err != nil {
		logger.Log.WithField("execPath", execPath).Debug("Found executable path")
		return fmt.Errorf("could not determine executable path: %w", err)
	}

	var daemonPort = "2388"

	configPort := app.Config.GetInt("port")
	if configPort > 0 {
		app.CacheTools.WriteStringToFile("wallhaven/daemon_port", strconv.Itoa(configPort))
		daemonPort = strconv.Itoa(configPort)
	} else {
		cachePort, err := app.CacheTools.ReadLineFromFile("wallhaven/daemon_port", 1)
		if err != nil {
			daemonPort = "2388"
			app.CacheTools.WriteStringToFile("wallhaven/daemon_port", strconv.Itoa(configPort))
		}
		if cachePort != "" {
			daemonPort = cachePort
		}
	}

	cmd := exec.Command(execPath, "-startdaemon", "-port", daemonPort)
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

	os.Exit(0)

	return nil

}

func KillDaemon(app *appcontext.AppContext) error {

	var port string
	port, err := app.CacheTools.ReadLineFromFile("wallhaven/daemon_port", 1)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	killStr := fmt.Sprintf("http://localhost:%s/kill", port)

	_, err = http.Get(killStr)
	if err != nil {
		logger.Log.Debug("Could not kill daemon with GET request")
	}

	logger.Log.Info("Daemon process stopped")

	return nil

}
