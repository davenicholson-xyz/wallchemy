package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type WallpaperDownloader struct {
	cacheTools   *CacheTools
	httpClient   *http.Client
	outputDir    string // Custom output directory
	useCustomDir bool   // Whether to use custom directory or cache
}

func NewWallpaperDownloader(cacheTools *CacheTools, outputDir string) *WallpaperDownloader {
	wd := &WallpaperDownloader{
		cacheTools: cacheTools,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// If output directory is specified, use it instead of cache
	if outputDir != "" {
		// Expand and clean the path
		expandedDir := expandPath(outputDir)

		// Convert to absolute path
		absDir, err := filepath.Abs(expandedDir)
		if err != nil {
			Logger.WithError(err).WithField("output_dir", outputDir).Warn("Failed to get absolute path, using as-is")
			absDir = expandedDir
		}

		wd.outputDir = absDir
		wd.useCustomDir = true

		// Ensure the custom output directory exists
		if err := os.MkdirAll(absDir, 0755); err != nil {
			Logger.WithError(err).WithField("output_dir", absDir).Warn("Failed to create custom output directory, falling back to cache")
			wd.useCustomDir = false
		} else {
			Logger.WithField("output_dir", absDir).Info("Using custom output directory")
		}
	}

	return wd
}

type DownloadResult struct {
	FilePath    string
	WasCached   bool
	FileSize    int64
	OriginalURL string
}

func (wd *WallpaperDownloader) GetOrDownloadWallpaper(wallpaperURL string) (*DownloadResult, error) {

	filename, err := getFilenameFromURL(wallpaperURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract filename from URL: %w", err)
	}

	// Determine where to store/look for the file
	var filePath string

	if wd.useCustomDir {
		// Use custom output directory exclusively
		filePath = filepath.Join(wd.outputDir, filename)
	} else {
		// Use cache directory
		cacheFileName := fmt.Sprintf("wallpapers/%s", filename)
		filePath = wd.cacheTools.Join(cacheFileName)
	}

	// Check if file already exists
	if info, err := os.Stat(filePath); err == nil {
		Logger.WithFields(logrus.Fields{
			"filename":  filename,
			"file_path": filePath,
			"file_size": info.Size(),
		}).Info("Using existing wallpaper")

		return &DownloadResult{
			FilePath:    filePath,
			WasCached:   true,
			FileSize:    info.Size(),
			OriginalURL: wallpaperURL,
		}, nil
	}

	Logger.WithFields(logrus.Fields{
		"filename": filename,
		"url":      wallpaperURL,
	}).Info("Downloading wallpaper")

	result, err := wd.downloadWallpaper(wallpaperURL, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to download wallpaper: %w", err)
	}

	return result, nil
}

func (wd *WallpaperDownloader) downloadWallpaper(wallpaperURL, filename string) (*DownloadResult, error) {

	resp, err := wd.httpClient.Get(wallpaperURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var finalPath string

	if wd.useCustomDir {
		// Save to custom directory only
		finalPath = filepath.Join(wd.outputDir, filename)
		if err := os.WriteFile(finalPath, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to save wallpaper to custom directory: %w", err)
		}

		Logger.WithFields(logrus.Fields{
			"file_path": finalPath,
			"file_size": len(data),
			"url":       wallpaperURL,
		}).Info("Successfully downloaded wallpaper to custom directory")

	} else {
		// Save to cache directory only
		cacheFileName := fmt.Sprintf("wallpapers/%s", filename)
		if err := wd.cacheTools.WriteFileBytes(cacheFileName, data); err != nil {
			return nil, fmt.Errorf("failed to save wallpaper to cache: %w", err)
		}

		finalPath = wd.cacheTools.Join(cacheFileName)

		Logger.WithFields(logrus.Fields{
			"file_path": finalPath,
			"file_size": len(data),
			"url":       wallpaperURL,
		}).Info("Successfully downloaded and cached wallpaper")
	}

	return &DownloadResult{
		FilePath:    finalPath,
		WasCached:   false,
		FileSize:    int64(len(data)),
		OriginalURL: wallpaperURL,
	}, nil
}

func getFilenameFromURL(wallpaperURL string) (string, error) {
	parsedURL, err := url.Parse(wallpaperURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	filename := filepath.Base(parsedURL.Path)
	if filename == "" || filename == "." || filename == "/" {
		return "", fmt.Errorf("no filename found in URL: %s", wallpaperURL)
	}

	return filename, nil
}

// expandPath handles tilde expansion and path cleaning
func expandPath(path string) string {
	// Handle tilde expansion
	if len(path) >= 2 && path[:2] == "~/" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(homeDir, path[2:])
		}
	}

	// Clean the path (handles ./ and ../ references)
	return filepath.Clean(path)
}
