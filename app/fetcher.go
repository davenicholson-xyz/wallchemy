package app

import (
	"fmt"
	"math/rand"
	"strings"

	wapi "github.com/davenicholson-xyz/go-wallhaven/wallhavenapi"
	"github.com/sirupsen/logrus"
)

type WallpaperFetcher struct {
	cacheTools *CacheTools
}

func NewWallpaperFetcher(cacheTools *CacheTools) *WallpaperFetcher {
	return &WallpaperFetcher{
		cacheTools: cacheTools,
	}
}

type FetchConfig struct {
	Query     *wapi.Query
	MaxPages  int
	CacheFile string
}

type FetchResult struct {
	Wallpapers      []string
	SuccessfulPages int
	RequestedPages  int
	SelectedRandom  string
}

type pageResult struct {
	page       int
	wallpapers []string
	err        error
}

func (wf *WallpaperFetcher) FetchAndCacheWallpapers(config FetchConfig) (*FetchResult, error) {
	Logger.WithFields(logrus.Fields{
		"max_pages":  config.MaxPages,
		"cache_file": config.CacheFile,
	}).Info("Fetching wallpapers from API")

	result, err := wf.fetchConcurrently(config)
	if err != nil {
		return nil, err
	}

	if len(result.Wallpapers) > 0 {
		cacheFileName := fmt.Sprintf("%s_history", config.CacheFile)
		if err := wf.saveToCache(cacheFileName, result.Wallpapers); err != nil {
			Logger.WithError(err).Warn("Failed to save wallpapers to cache")
		} else {
			Logger.WithFields(logrus.Fields{
				"cache_file":      cacheFileName,
				"wallpaper_count": len(result.Wallpapers),
			}).Info("Successfully cached wallpapers")
		}
	}

	if len(result.Wallpapers) > 0 {
		result.SelectedRandom = result.Wallpapers[rand.Intn(len(result.Wallpapers))]
		Logger.WithFields(logrus.Fields{
			"total_available":    len(result.Wallpapers),
			"selected_wallpaper": result.SelectedRandom,
		}).Info("Selected random wallpaper from fetched results")
	}

	return result, nil
}

func (wf *WallpaperFetcher) fetchConcurrently(config FetchConfig) (*FetchResult, error) {

	Logger.WithField("page", 1).Debug("Fetching first page to determine total pages")

	firstResp, err := config.Query.Page(1)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch first page: %w", err)
	}

	totalPages := firstResp.Meta.LastPage
	if totalPages == 0 {
		totalPages = 1
	}

	pagesToFetch := min(totalPages, config.MaxPages)

	Logger.WithFields(logrus.Fields{
		"total_pages":    totalPages,
		"pages_to_fetch": pagesToFetch,
		"requested_max":  config.MaxPages,
	}).Info("Determined page count for fetching")

	firstPageWallpapers := make([]string, 0, len(firstResp.Wallpapers))
	for _, wp := range firstResp.Wallpapers {
		firstPageWallpapers = append(firstPageWallpapers, wp.Path)
	}

	Logger.WithFields(logrus.Fields{
		"page":            1,
		"page_wallpapers": len(firstPageWallpapers),
	}).Debug("Successfully fetched first page")

	if len(firstPageWallpapers) == 0 {
		Logger.Info("No wallpapers found")
		return &FetchResult{}, nil
	}

	if pagesToFetch == 1 {
		Logger.Info("Only one page needed, returning first page results")
		return &FetchResult{
			Wallpapers:      firstPageWallpapers,
			SuccessfulPages: 1,
			RequestedPages:  pagesToFetch,
		}, nil
	}

	remainingPages := pagesToFetch - 1
	results := make(chan pageResult, remainingPages)

	for page := 2; page <= pagesToFetch; page++ {
		go func(p int) {
			Logger.WithField("page", p).Debug("Fetching page")

			resp, err := config.Query.Page(p)
			if err != nil {
				Logger.WithFields(logrus.Fields{
					"page":  p,
					"error": err,
				}).Warn("Failed to fetch page")
				results <- pageResult{page: p, err: err}
				return
			}

			pageWallpapers := make([]string, 0, len(resp.Wallpapers))
			for _, wp := range resp.Wallpapers {
				pageWallpapers = append(pageWallpapers, wp.Path)
			}

			Logger.WithFields(logrus.Fields{
				"page":            p,
				"page_wallpapers": len(pageWallpapers),
			}).Debug("Successfully fetched page")

			results <- pageResult{page: p, wallpapers: pageWallpapers}
		}(page)
	}

	allWallpapers := firstPageWallpapers
	successfulPages := 1

	for range remainingPages {
		result := <-results

		if result.err != nil {
			continue
		}

		if len(result.wallpapers) == 0 {
			Logger.WithField("page", result.page).Debug("No wallpapers found on page")
			continue
		}

		allWallpapers = append(allWallpapers, result.wallpapers...)
		successfulPages++

		Logger.WithFields(logrus.Fields{
			"page":            result.page,
			"page_wallpapers": len(result.wallpapers),
			"total_so_far":    len(allWallpapers),
		}).Debug("Collected wallpapers from page")
	}

	if len(allWallpapers) == 0 {
		return nil, fmt.Errorf("no wallpapers found")
	}

	Logger.WithFields(logrus.Fields{
		"total_wallpapers": len(allWallpapers),
		"successful_pages": successfulPages,
		"requested_pages":  pagesToFetch,
		"actual_total":     totalPages,
	}).Info("Successfully fetched wallpapers concurrently")

	return &FetchResult{
		Wallpapers:      allWallpapers,
		SuccessfulPages: successfulPages,
		RequestedPages:  pagesToFetch,
	}, nil
}

func (wf *WallpaperFetcher) saveToCache(cacheFile string, wallpapers []string) error {
	if err := wf.cacheTools.WriteStringToFile(cacheFile, strings.Join(wallpapers, "\n")); err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	Logger.WithFields(logrus.Fields{
		"cache_file":      cacheFile,
		"wallpaper_count": len(wallpapers),
	}).Debug("Saved wallpapers to cache")

	return nil
}
