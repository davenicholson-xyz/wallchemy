package handlers

import (
	"fmt"

	"github.com/davenicholson-xyz/wallchemy/app"
	"github.com/sirupsen/logrus"
)

func handleSingle(appCtx *app.AppContext) (string, error) {

	wp, err := appCtx.WallhavenAPI.Wallpaper(appCtx.Config.ID)
	if err != nil {
		return "", err
	}

	downloader := app.NewWallpaperDownloader(appCtx.CacheTools, appCtx.Config.Output)
	downloadResult, err := downloader.GetOrDownloadWallpaper(wp.Path)
	if err != nil {
		return "", fmt.Errorf("failed to download wallpaper: %w", err)
	}

	app.Logger.WithFields(logrus.Fields{
		"file_path":  downloadResult.FilePath,
		"was_cached": downloadResult.WasCached,
		"file_size":  downloadResult.FileSize,
	}).Info("Single wallpaper ready")

	if err := app.ApplyWallpaper(downloadResult.FilePath, appCtx); err != nil {
		return "", err
	}

	return downloadResult.FilePath, nil
}
