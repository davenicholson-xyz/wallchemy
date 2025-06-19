package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"
)

type CacheTools struct {
	cacheDir string
}

type FileInfo struct {
	Name        string
	Size        int64
	ModTime     time.Time
	IsDir       bool
	Permissions os.FileMode
	FullPath    string
}

func NewCacheTools(path string) (*CacheTools, error) {
	ct := &CacheTools{}
	systemCache := getSystemCache()
	err := ct.SetCacheDir(filepath.Join(systemCache, path))
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
	return filepath.Join(pathSlice...)
}

func (ct *CacheTools) Clear() error {
	cacheDir := ct.GetCacheDirectory()
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirPath := filepath.Join(cacheDir, entry.Name())
		err := os.RemoveAll(dirPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ct *CacheTools) Clean() error {
	entries, err := os.ReadDir(ct.cacheDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(ct.cacheDir, entry.Name())
		if err := os.RemoveAll(fullPath); err != nil {
			return err
		}
	}

	return nil
}

func (ct *CacheTools) ClearAll() error {
	return os.RemoveAll(ct.cacheDir)
}

func (ct *CacheTools) WriteStringToFile(filename string, str string) error {
	fullPath := ct.Join(filename)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directories: %w", err)
	}
	if err := os.WriteFile(fullPath, []byte(str), 0600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}
	return nil
}

func (ct *CacheTools) AppendStringToFile(filename string, str string) error {
	fullPath := ct.Join(filename)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directories: %w", err)
	}

	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file for appending: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(str)
	if err != nil {
		return fmt.Errorf("failed to append to file: %w", err)
	}
	return nil
}

func (ct *CacheTools) ReadLineFromFile(filename string, line int) (string, error) {
	if line < 1 {
		return "", fmt.Errorf("line number must be >= 1")
	}

	fullPath := ct.Join(filename)
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
	return "", fmt.Errorf("line %d not found", line)
}

func (ct *CacheTools) ReadLines(filename string) ([]string, error) {
	fullPath := ct.Join(filename)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func (ct *CacheTools) WriteLines(filename string, lines []string) error {
	return ct.WriteStringToFile(filename, strings.Join(lines, "\n"))
}

func (ct *CacheTools) ReplaceLineInFile(filename string, lineNum int, newContent string) error {
	lines, err := ct.ReadLines(filename)
	if err != nil {
		return err
	}

	if lineNum < 1 || lineNum > len(lines) {
		return fmt.Errorf("line number %d is out of range (1-%d)", lineNum, len(lines))
	}

	lines[lineNum-1] = newContent
	return ct.WriteLines(filename, lines)
}

func (ct *CacheTools) InsertLineInFile(filename string, lineNum int, content string) error {
	lines, err := ct.ReadLines(filename)
	if err != nil {
		return err
	}

	if lineNum < 1 || lineNum > len(lines)+1 {
		return fmt.Errorf("line number %d is out of range (1-%d)", lineNum, len(lines)+1)
	}

	newLines := make([]string, len(lines)+1)
	copy(newLines[:lineNum-1], lines[:lineNum-1])
	newLines[lineNum-1] = content
	copy(newLines[lineNum:], lines[lineNum-1:])

	return ct.WriteLines(filename, newLines)
}

func (ct *CacheTools) DeleteLineFromFile(filename string, lineNum int) error {
	lines, err := ct.ReadLines(filename)
	if err != nil {
		return err
	}

	if lineNum < 1 || lineNum > len(lines) {
		return fmt.Errorf("line number %d is out of range (1-%d)", lineNum, len(lines))
	}

	newLines := slices.Delete(lines, lineNum-1, lineNum)
	return ct.WriteLines(filename, newLines)
}

func (ct *CacheTools) DeleteFile(filename string) error {
	fullPath := ct.Join(filename)
	if err := os.Remove(fullPath); err != nil {
		return err
	}
	return nil
}

func (ct *CacheTools) ReadFile(filename string) (string, error) {
	fullPath := ct.Join(filename)
	b, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (ct *CacheTools) ReadFileBytes(filename string) ([]byte, error) {
	fullPath := ct.Join(filename)
	return os.ReadFile(fullPath)
}

func (ct *CacheTools) WriteFileBytes(filename string, data []byte) error {
	fullPath := ct.Join(filename)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directories: %w", err)
	}
	if err := os.WriteFile(fullPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}
	return nil
}

func (ct *CacheTools) CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstPath := ct.Join(dst)

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directories: %w", err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (ct *CacheTools) FileExists(filename string) bool {
	fullPath := ct.Join(filename)
	_, err := os.Stat(fullPath)
	return err == nil
}

func (ct *CacheTools) GetFileInfo(filename string) (*FileInfo, error) {
	fullPath := ct.Join(filename)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:        stat.Name(),
		Size:        stat.Size(),
		ModTime:     stat.ModTime(),
		IsDir:       stat.IsDir(),
		Permissions: stat.Mode(),
		FullPath:    fullPath,
	}, nil
}

func (ct *CacheTools) GetFileModTime(filename string) (time.Time, error) {
	info, err := ct.GetFileInfo(filename)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime, nil
}

func (ct *CacheTools) GetFileSize(filename string) (int64, error) {
	info, err := ct.GetFileInfo(filename)
	if err != nil {
		return 0, err
	}
	return info.Size, nil
}

func (ct *CacheTools) ListFiles() ([]*FileInfo, error) {
	entries, err := os.ReadDir(ct.cacheDir)
	if err != nil {
		return nil, err
	}

	var files []*FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileInfo := &FileInfo{
			Name:        info.Name(),
			Size:        info.Size(),
			ModTime:     info.ModTime(),
			IsDir:       info.IsDir(),
			Permissions: info.Mode(),
			FullPath:    ct.Join(info.Name()),
		}
		files = append(files, fileInfo)
	}

	return files, nil
}

func (ct *CacheTools) ListFilesRecursive() ([]*FileInfo, error) {
	var files []*FileInfo

	err := filepath.Walk(ct.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(ct.cacheDir, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		fileInfo := &FileInfo{
			Name:        filepath.Base(path),
			Size:        info.Size(),
			ModTime:     info.ModTime(),
			IsDir:       info.IsDir(),
			Permissions: info.Mode(),
			FullPath:    path,
		}
		files = append(files, fileInfo)

		return nil
	})

	return files, err
}

func (ct *CacheTools) CreateSubdir(dirname string) error {
	fullPath := ct.Join(dirname)
	return os.MkdirAll(fullPath, 0755)
}

func (ct *CacheTools) RemoveSubdir(dirname string) error {
	fullPath := ct.Join(dirname)
	return os.RemoveAll(fullPath)
}

func getSystemCache() string {
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("LOCALAPPDATA"); appData != "" {
			return appData
		}
		return filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	case "linux", "freebsd", "openbsd", "netbsd":
		if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
			return xdg
		}
		return filepath.Join(os.Getenv("HOME"), ".cache")
	default:
		return "./cache"
	}
}
