package cachetools

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type CacheTools struct {
	cacheDir string
}

func New(path string) (*CacheTools, error) {
	ct := &CacheTools{}
	system_cache := getSystemCache()
	err := ct.SetCacheDir(filepath.Join(system_cache, path))
	if err != nil {
		return nil, err
	}
	return ct, nil
}

func (ct *CacheTools) SetCacheDir(path string) error {
	ct.cacheDir = path
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create cache directories: %w", err)
	}
	return nil
}

func (ct *CacheTools) GetCacheDirectory() string {
	return ct.cacheDir
}

func (ct *CacheTools) Join(path ...string) string {
	pathSlice := append([]string{ct.GetCacheDirectory()}, path...)
	fullPath := filepath.Join(pathSlice...)
	return fullPath
}

func (ct *CacheTools) Clear() error {
	err := os.RemoveAll(ct.GetCacheDirectory())
	if err != nil {
		return err
	}
	return nil
}

func (ct *CacheTools) WriteStringToFile(filename string, str string) error {
	fullPath := filepath.Join(ct.GetCacheDirectory(), filename)
	if err := os.WriteFile(fullPath, []byte(str), 0600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}
	return nil
}

func (ct *CacheTools) ReadLineFromFile(filename string, line int) (string, error) {
	fullPath := filepath.Join(ct.GetCacheDirectory(), filename)
	file, err := os.Open(fullPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 1

	for scanner.Scan() {
		if currentLine == line {
			return scanner.Text(), nil
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func (ct *CacheTools) ReadFile(filename string) (string, error) {
	fullPath := filepath.Join(ct.GetCacheDirectory(), filename)
	b, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(b), nil

}

func getSystemCache() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("LocalAppData"))
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	case "linux", "freebsd", "openbsd", "netbsd":
		if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
			return filepath.Join(xdg, "myapp")
		}
		return filepath.Join(os.Getenv("HOME"), ".cache")
	default:
		return "./cache"
	}
}
