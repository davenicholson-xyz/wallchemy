package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/davenicholson-xyz/wallchemy/app"
	"github.com/sirupsen/logrus"
)

type UsersCollection struct {
	ID     int    `json:"id"`
	Label  string `json:"label"`
	Views  int    `json:"views"`
	Public int    `json:"public"`
	Count  int    `json:"count"`
}

type Collections struct {
	Data []UsersCollection `json:"data"`
}

func handleCollections(appCtx *app.AppContext) (string, error) {

	if appCtx.Config.APIKey == "" {
		return "", fmt.Errorf("API key is required to get collections")
	}

	collections, err := getOrFetchCollections(appCtx)
	if err != nil {
		return "", fmt.Errorf("failed to get collections: %w", err)
	}

	var availableCollections []string
	for _, collection := range collections.Data {
		availableCollections = append(availableCollections, collection.Label)
	}

	return fmt.Sprintf("%s", strings.Join(availableCollections, "\n")), nil

}

func handleCollection(appCtx *app.AppContext) (string, error) {
	var selectedWallpaper string

	if err := globalFilters(appCtx); err != nil {
		return "", err
	}

	collectionName := appCtx.Config.Collection
	if collectionName == "" {
		return "", fmt.Errorf("no collection name specified")
	}

	if appCtx.Config.APIKey == "" {
		return "", fmt.Errorf("API key is required to get collections")
	}

	if appCtx.Config.Username == "" {
		return "", fmt.Errorf("API key is required to get collections")
	}

	collections, err := getOrFetchCollections(appCtx)
	if err != nil {
		return "", fmt.Errorf("failed to get collections: %w", err)
	}

	var targetCollection *UsersCollection
	for _, collection := range collections.Data {
		if strings.EqualFold(collection.Label, collectionName) {
			targetCollection = &collection
			break
		}
	}

	if targetCollection == nil {
		var availableCollections []string
		for _, collection := range collections.Data {
			availableCollections = append(availableCollections, collection.Label)
		}
		return "", fmt.Errorf("collection '%s' not found. Available collections: %s",
			collectionName, strings.Join(availableCollections, ", "))
	}

	app.Logger.WithFields(logrus.Fields{
		"collection_id":   targetCollection.ID,
		"collection_name": targetCollection.Label,
		"wallpaper_count": targetCollection.Count,
	}).Info("Found target collection")

	cacheFileName := fmt.Sprintf("collection_%d", targetCollection.ID)
	cachedWallpaper, err := app.CheckCacheByName(cacheFileName, appCtx)
	if err != nil {
		app.Logger.WithError(err).Warn("Error checking collection cache, will fetch from API")
	} else if cachedWallpaper != "" {
		selectedWallpaper = cachedWallpaper
		app.Logger.WithField("wallpaper", selectedWallpaper).Info("Using cached collection wallpaper")
	}

	if selectedWallpaper == "" {

		app.Logger.Info("No valid cache found, fetching from API")

		query := appCtx.WallhavenAPI.Collection(appCtx.Config.Username, targetCollection.ID)

		fetcher := app.NewWallpaperFetcher(appCtx.CacheTools)
		fetchConfig := app.FetchConfig{
			Query:     query,
			MaxPages:  appCtx.Config.MaxPages,
			CacheFile: cacheFileName,
		}

		result, err := fetcher.FetchAndCacheWallpapers(fetchConfig)
		if err != nil {
			return "", err
		}

		selectedWallpaper = result.SelectedRandom
	}

	if selectedWallpaper == "" {
		appCtx.CacheTools.DeleteFile(fmt.Sprintf("%s_history", cacheFileName))
		return "", fmt.Errorf("No wallpapers found")
	}

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

	return downloadResult.FilePath, nil

}

func getOrFetchCollections(appCtx *app.AppContext) (*Collections, error) {
	collectionsFile := "collections"

	if appCtx.CacheTools.FileExists(collectionsFile) {
		collections, err := readCollectionsFromCache(appCtx, collectionsFile)
		if err != nil {
			app.Logger.WithError(err).Warn("Error reading collections from cache, will fetch fresh")
		} else {
			app.Logger.WithField("collections_count", len(collections.Data)).Info("Using cached collections")
			return collections, nil
		}
	}

	app.Logger.Info("Collections cache not found, fetching from API")

	collections, err := fetchCollectionsFromAPI(appCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collections from API: %w", err)
	}

	if err := saveCollectionsToCache(appCtx, collectionsFile, collections); err != nil {
		app.Logger.WithError(err).Warn("Failed to save collections to cache")
	}

	app.Logger.WithField("collections_count", len(collections.Data)).Info("Fetched and cached collections")
	return collections, nil
}

func fetchCollectionsFromAPI(appCtx *app.AppContext) (*Collections, error) {
	collectionsData, err := appCtx.WallhavenAPI.MyCollections()
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	var collections Collections
	collections.Data = make([]UsersCollection, len(collectionsData))

	for i, col := range collectionsData {
		collections.Data[i] = UsersCollection{
			ID:     col.ID,
			Label:  col.Label,
			Views:  col.Views,
			Public: col.Public,
			Count:  col.Count,
		}
	}

	return &collections, nil
}

func saveCollectionsToCache(appCtx *app.AppContext, filename string, collections *Collections) error {
	data, err := json.MarshalIndent(collections, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal collections: %w", err)
	}

	if err := appCtx.CacheTools.WriteStringToFile(filename, string(data)); err != nil {
		return fmt.Errorf("failed to write collections to cache: %w", err)
	}

	app.Logger.WithField("file", filename).Debug("Saved collections to cache")
	return nil
}

func readCollectionsFromCache(appCtx *app.AppContext, filename string) (*Collections, error) {
	data, err := appCtx.CacheTools.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read collections from cache: %w", err)
	}

	var collections Collections
	if err := json.Unmarshal([]byte(data), &collections); err != nil {
		return nil, fmt.Errorf("failed to unmarshal collections: %w", err)
	}

	return &collections, nil
}
