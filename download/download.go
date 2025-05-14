package download

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
)

func FetchJson(url string) ([]byte, error) {
	slog.Info("Fetching:" + url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	return body, nil
}

func GenerateSeed(length int) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	b := make([]byte, length)
	for i := range length {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func DownloadImage(url string, output string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
