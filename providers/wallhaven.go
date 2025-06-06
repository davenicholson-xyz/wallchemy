package providers

import (
	"fmt"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/logger"
	wh "github.com/davenicholson-xyz/wallchemy/providers/wallhaven"
)

type WallhavenProvider struct{}

func (w *WallhavenProvider) Name() string {
	return "wallhaven"
}

func (w *WallhavenProvider) ParseArgs(app *appcontext.AppContext) (string, error) {

	if app.Config.GetBool("daemon") {
		wh.KillDaemon(app)
		wh.LaunchDaemon(app)
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

	if app.Config.GetBool("last") {
		return wh.LastRandomQuery(app), nil
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
		logger.Log.WithField("selected", selection).Info("Selected wallpaper")
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
