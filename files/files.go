package files

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func IsFullPath(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	return absPath == filepath.Clean(path)
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		slog.Debug("Path does not exist: " + path)
		return false
	}
	slog.Debug("Path exists: " + path)
	return true
}

func GetUserConfigDir() (string, bool) {
	var configpath string
	if runtime.GOOS == "windows" {
		configpath = filepath.Join(os.Getenv("APPDATA"), "wallchemy")
		if PathExists(configpath) {
			return configpath, true
		} else {
			return configpath, false
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		configpath = filepath.Join("/tmp")
	}
	configpath = filepath.Join(home, ".config", "wallchemy")

	if PathExists(configpath) {
		return configpath, true
	} else {
		return configpath, false
	}
}

func DefaultConfigFilepath() (string, bool) {
	cfg_dir, _ := GetUserConfigDir()
	cfg_path := filepath.Join(cfg_dir, "config.yml")
	exists := PathExists(cfg_path)
	return cfg_path, exists
}

func IsFileFresh(filepath string, expirySeconds int) bool {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	modifiedTime := fileInfo.ModTime()
	expiryDuration := time.Duration(expirySeconds) * time.Second
	expiryTime := modifiedTime.Add(expiryDuration)

	return time.Now().Before(expiryTime)
}

func GetRandomLine(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("Failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	if len(nonEmptyLines) == 0 {
		return "", fmt.Errorf("No wallpapers found")
	}

	return nonEmptyLines[rand.Intn(len(nonEmptyLines))], nil
}
