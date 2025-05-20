package download

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/davenicholson-xyz/wallchemy/logger"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
)

func FetchJson(url string) ([]byte, error) {
	logger.Log.WithField("URL", url).Debug("Fetching url")

	resp, err := http.Get(url)
	if err != nil {
		logger.Log.WithError(err)
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.WithError(err)
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Error:404")
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("Error:401")
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
	logger.Log.WithField("Image url:", url).Info("Downloading image")
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
