package app

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	wapi "github.com/davenicholson-xyz/go-wallhaven/wallhavenapi"
	"github.com/sirupsen/logrus"
)

type CacheManager struct {
	cacheTools *CacheTools
}

func NewCacheManager(cacheTools *CacheTools) *CacheManager {
	return &CacheManager{
		cacheTools: cacheTools,
	}
}

func (cm *CacheManager) getFileNames(sort wapi.SortingType) (queryFile, cacheFile string) {
	sortStr := string(sort)
	return fmt.Sprintf("%s_query", sortStr), fmt.Sprintf("%s_history", sortStr)
}

func (cm *CacheManager) buildCurrentURL(wallhavenAPI *wapi.WallhavenAPI, query string, sort wapi.SortingType) (string, error) {
	qUrl := wallhavenAPI.Search(query).Sort(sort)
	return RemoveURLParams(qUrl.Raw(), "apikey", "seed")
}

func (cm *CacheManager) hasQueryChanged(queryFile, currentURL string) (bool, error) {
	if !cm.cacheTools.FileExists(queryFile) {
		Logger.Debug("Query file doesn't exist")
		return true, nil
	}

	prevURL, err := cm.cacheTools.ReadFile(queryFile)
	if err != nil {
		return true, fmt.Errorf("failed to read previous query: %w", err)
	}

	prevURL = strings.TrimSpace(prevURL)
	changed := prevURL != currentURL

	if changed {
		Logger.WithFields(logrus.Fields{
			"previous_url": prevURL,
			"current_url":  currentURL,
		}).Debug("Query has changed")
	} else {
		Logger.WithField("query", currentURL).Debug("Query unchanged")
	}

	return changed, nil
}

func (cm *CacheManager) updateQueryFile(queryFile, currentURL string) error {
	if err := cm.cacheTools.WriteStringToFile(queryFile, currentURL); err != nil {
		return fmt.Errorf("failed to update query file: %w", err)
	}
	Logger.WithField("query_file", queryFile).Debug("Updated query file")
	return nil
}

func (cm *CacheManager) isCacheStale(cacheFile string, expirySeconds int) (bool, error) {
	if !cm.cacheTools.FileExists(cacheFile) {
		Logger.WithField("cache_file", cacheFile).Debug("Cache file doesn't exist")
		return true, nil
	}

	modTime, err := cm.cacheTools.GetFileModTime(cacheFile)
	if err != nil {
		return true, fmt.Errorf("failed to get cache file modification time: %w", err)
	}

	expiryDuration := time.Duration(expirySeconds) * time.Second
	age := time.Since(modTime)
	isStale := age > expiryDuration

	Logger.WithFields(logrus.Fields{
		"cache_file":    cacheFile,
		"file_age":      age.String(),
		"expiry_limit":  expiryDuration.String(),
		"is_stale":      isStale,
		"last_modified": modTime.Format(time.RFC3339),
	}).Debug("Cache staleness check")

	if isStale {
		Logger.WithFields(logrus.Fields{
			"cache_file": cacheFile,
			"age":        age.String(),
			"limit":      expiryDuration.String(),
		}).Info("Cache file is stale and will be ignored")
	}

	return isStale, nil
}

func (cm *CacheManager) getRandomWallpaperFromCache(cacheFile string, expirySeconds int) (string, error) {

	stale, err := cm.isCacheStale(cacheFile, expirySeconds)
	if err != nil {
		return "", err
	}
	if stale {
		return "", nil
	}

	lines, err := cm.cacheTools.ReadLines(cacheFile)
	if err != nil {
		return "", fmt.Errorf("failed to read cache file: %w", err)
	}

	if len(lines) == 0 {
		Logger.WithField("cache_file", cacheFile).Debug("Cache file is empty")
		return "", nil
	}

	var validLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			validLines = append(validLines, strings.TrimSpace(line))
		}
	}

	if len(validLines) == 0 {
		Logger.WithField("cache_file", cacheFile).Debug("No valid wallpapers in cache file")
		return "", nil
	}

	selectedWallpaper := validLines[rand.Intn(len(validLines))]

	Logger.WithFields(logrus.Fields{
		"cache_file":         cacheFile,
		"total_wallpapers":   len(validLines),
		"selected_wallpaper": selectedWallpaper,
	}).Debug("Selected wallpaper from cache")

	return selectedWallpaper, nil
}

func (cm *CacheManager) CheckForCache(wallhavenAPI *wapi.WallhavenAPI, query string, sort wapi.SortingType, expirySeconds int) (string, error) {
	queryFile, cacheFile := cm.getFileNames(sort)

	currentURL, err := cm.buildCurrentURL(wallhavenAPI, query, sort)
	if err != nil {
		return "", fmt.Errorf("failed to build current URL: %w", err)
	}

	queryChanged, err := cm.hasQueryChanged(queryFile, currentURL)
	if err != nil {
		return "", fmt.Errorf("failed to check query changes: %w", err)
	}

	if queryChanged {
		if err := cm.updateQueryFile(queryFile, currentURL); err != nil {
			return "", err
		}

		return "", nil
	}

	return cm.getRandomWallpaperFromCache(cacheFile, expirySeconds)
}

func CheckForCache(appCtx *AppContext) (string, error) {
	cacheManager := NewCacheManager(appCtx.CacheTools)

	query := appCtx.Config.Random
	sort := SortMethod(appCtx)
	expirySeconds := appCtx.Config.Expiry

	return cacheManager.CheckForCache(appCtx.WallhavenAPI, query, sort, expirySeconds)
}

func CheckCacheByName(fileName string, appCtx *AppContext) (string, error) {
	cacheManager := NewCacheManager(appCtx.CacheTools)

	queryFile := fmt.Sprintf("%s_query", fileName)
	cacheFile := fmt.Sprintf("%s_history", fileName)

	currentURL, err := cacheManager.buildCurrentURL(
		appCtx.WallhavenAPI,
		appCtx.Config.Random,
		SortMethod(appCtx),
	)
	if err != nil {
		return "", fmt.Errorf("failed to build current URL: %w", err)
	}

	queryChanged, err := cacheManager.hasQueryChanged(queryFile, currentURL)
	if err != nil {
		Logger.WithError(err).Warn("Error checking query changes")
		return "", nil
	}

	if queryChanged {
		Logger.WithField("file", fileName).Debug("Query changed, cache invalid")
		if err := cacheManager.updateQueryFile(queryFile, currentURL); err != nil {
			Logger.WithError(err).Warn("Failed to update query file")
		}
		return "", nil
	}

	cachedWallpaper, err := cacheManager.getRandomWallpaperFromCache(cacheFile, appCtx.Config.Expiry)
	if err != nil {
		Logger.WithError(err).Warn("Error reading from cache")
		return "", nil
	}

	if cachedWallpaper != "" {
		Logger.WithFields(logrus.Fields{
			"file":      fileName,
			"wallpaper": cachedWallpaper,
		}).Info("Found valid cached wallpaper")
	} else {
		Logger.WithField("file", fileName).Debug("No valid cache found")
		return "", nil
	}

	return cachedWallpaper, nil
}
