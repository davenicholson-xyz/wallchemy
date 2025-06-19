package app

import (
	"os"
	"os/exec"

	"github.com/davenicholson-xyz/go-setwallpaper/wallpaper"
)

func ApplyWallpaper(wp string, appCtx *AppContext) error {
	if appCtx.Config.Stealth {
		Logger.Info("Stealth mode - skipping setting wallpaper")
	} else {
		script := appCtx.Config.Script
		if script != "" {
			Logger.WithField("script", script).Debug("Setting wallpaper with script")
			cmd := exec.Command("sh", "-c", script+" "+wp)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return err
			}
		} else {
			Logger.Debug("Setting wallpaper with built-in")
			wallpaper.Set(wp)
		}

	}

	return nil
}
