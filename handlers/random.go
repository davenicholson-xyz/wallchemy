package handlers

import (
	"fmt"

	"github.com/davenicholson-xyz/wallchemy/app"
	"github.com/sirupsen/logrus"
)

func handleRandom(appCtx *app.AppContext) (string, error) {

	var selectedWallpaper string

	if err := globalFilters(appCtx); err != nil {
		return "", err
	}

	sortType := app.SortMethod(appCtx)
	fileName := string(sortType)

	cachedWallpaper, err := app.CheckCacheByName(fileName, appCtx)
	if err != nil {
		app.Logger.WithError(err).Warn("Error checking cache, will fetch from API")
	} else if cachedWallpaper != "" {
		app.Logger.WithField("wallpaper", cachedWallpaper).Info("Using cached wallpaper")
		selectedWallpaper = cachedWallpaper
	}

	if selectedWallpaper == "" {

		app.Logger.Info("No valid cache found, fetching from API")

		query := appCtx.WallhavenAPI.Search(appCtx.Config.Random).Sort(sortType)

		fetcher := app.NewWallpaperFetcher(appCtx.CacheTools)
		fetchConfig := app.FetchConfig{
			Query:     query,
			MaxPages:  appCtx.Config.MaxPages,
			CacheFile: fileName,
		}

		result, err := fetcher.FetchAndCacheWallpapers(fetchConfig)
		if err != nil {
			return "", err
		}

		selectedWallpaper = result.SelectedRandom
	}

	if selectedWallpaper == "" {
		appCtx.CacheTools.DeleteFile(fmt.Sprintf("%s_history", sortType))
		return "", fmt.Errorf("No wallpapers found")
	}

	app.Logger.WithField("wallpaper", selectedWallpaper).Info("Fetched and selected new wallpaper from API")

	downloader := app.NewWallpaperDownloader(appCtx.CacheTools, appCtx.Config.Output)
	downloadResult, err := downloader.GetOrDownloadWallpaper(selectedWallpaper)
	if err != nil {
		return "", fmt.Errorf("failed to download wallpaper: %w", err)
	}

	app.Logger.WithFields(logrus.Fields{
		"file_path":  downloadResult.FilePath,
		"was_cached": downloadResult.WasCached,
		"file_size":  downloadResult.FileSize,
	}).Info("Wallpaper ready")

	if err := app.ApplyWallpaper(downloadResult.FilePath, appCtx); err != nil {
		return "", err
	}

	appCtx.Current.Path = selectedWallpaper
	appCtx.Current.FilePath = downloadResult.FilePath
	SaveCurrent(appCtx)

	return downloadResult.FilePath, nil
}
