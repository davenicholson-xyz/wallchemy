package files

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/davenicholson-xyz/go-setwallpaper/wallpaper"
	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/download"
	"github.com/davenicholson-xyz/wallchemy/logger"
)

func ApplyWallpaper(wp string, app *appcontext.AppContext) (string, error) {

	filename := filepath.Base(wp)
	output := app.CacheTools.Join(app.Config.GetString("provider"), filename)

	if PathExists(output) {
		logger.Log.WithField("output", output).Debug("Already esists. Using cached")
	} else {
		logger.Log.WithField("output", output).Debug("Not in cache. Fetching")
		err := download.DownloadImage(wp, output)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
	}

	if app.Config.GetBool("stealth") {
		return output, nil
	}

	script := app.Config.GetString("script")
	if script != "" {
		logger.Log.WithField("script", script).Debug("Setting wallpaper with script")
		cmd := exec.Command("sh", "-c", script+" "+output)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return "", err
		}
	} else {
		logger.Log.Debug("Setting wallpaper with built-in")
		wallpaper.Set(output)
	}

	return output, nil
}
