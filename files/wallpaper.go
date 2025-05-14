package files

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/davenicholson-xyz/go-setwallpaper/wallpaper"
	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/download"
)

func ApplyWallpaper(wp string, app *appcontext.AppContext) (string, error) {

	filename := filepath.Base(wp)
	output := app.CacheTools.Join(app.Config.GetString("provider"), filename)

	if PathExists(output) {
		slog.Debug(output + " already exists. using cached version")
	} else {
		slog.Debug(output + " is new. fetching")
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
		cmd := exec.Command("sh", "-c", script+" "+output)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return "", err
		}
	} else {
		slog.Debug("Setting wallpaper with built in")
		wallpaper.Set(output)
	}

	return output, nil
}
