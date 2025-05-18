package providers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
	wh "github.com/davenicholson-xyz/wallchemy/providers/wallhaven"
)

type WallhavenProvider struct{}

func (w *WallhavenProvider) Name() string {
	return "wallhaven"
}

func isDaemonRunning(app *appcontext.AppContext) (bool, int) {
	data, err := app.CacheTools.ReadLineFromFile("daemon.pid", 1)
	if err != nil {
		return false, 0 // PID file doesn't exist
	}

	pid, err := strconv.Atoi(strings.TrimSpace(data))
	if err != nil {
		return false, 0 // Invalid PID
	}

	// Try to signal the process with 0 (doesn't kill it, just checks if it's there)
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, 0
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false, 0
	}

	return true, pid
}

func (w *WallhavenProvider) ParseArgs(app *appcontext.AppContext) (string, error) {

	if app.Config.GetBool("daemon") {
		running, pid := isDaemonRunning(app)
		if running {
			fmt.Printf("Daemon is already running with PID %d\n", pid)
			return "", nil
		}
		wh.LaunchDaemon()
	}

	if app.Config.GetBool("kill") {
		wh.KillDaemon(app)
	}

	if app.Config.GetBool("startdaemon") {
		wh.StartDaemon(app)
	}

	if app.Config.GetBool("file") {
		return wh.CurrentWallpaperInfo(app, 4), nil
	}

	if app.Config.GetBool("url") {
		return wh.CurrentWallpaperInfo(app, 3), nil
	}

	if app.Config.GetBool("link") {
		return wh.CurrentWallpaperInfo(app, 2), nil
	}

	var selection string
	var err error

	if app.Config.GetString("collection") != "" {
		selection, err = wh.SelectionFromCollection(app)
		if err != nil {
			return "", fmt.Errorf("Could not select from collection: %w", err)
		}
	}

	if app.Config.GetString("random") != "" || app.Config.GetBool("top") || app.Config.GetBool("hot") {
		selection, err = wh.SelectionFromRandom(app)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
	}

	if app.Config.GetString("id") != "" {
		selection, err = wh.SelectionFromID(app)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
	}

	if selection != "" {
		output, err := wh.SetSelectedWallpaper(selection, app)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		if output != "" {
			return output, nil
		}
	}

	return "", nil
}
