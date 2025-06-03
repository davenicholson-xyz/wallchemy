package wallhaven

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"sync"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/download"
	"github.com/davenicholson-xyz/wallchemy/files"
	"github.com/davenicholson-xyz/wallchemy/logger"
	"github.com/davenicholson-xyz/wallchemy/network"
)

type WallpaperData struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Path string `json:"path"`
}

type WallhavenRequest struct {
	Wallpapers []WallpaperData `json:"data"`
	Meta       WallhavenMeta   `json:"meta"`
}

type WallhavenMeta struct {
	LastPage int `json:"last_page"`
	Total    int `json:"total"`
}

func CurrentWallpaperInfo(app *appcontext.AppContext, line int) string {
	current := fmt.Sprintf("wallhaven/%s", "current")

	slog.Debug("output info from " + current)

	filename, err := app.CacheTools.ReadLineFromFile(current, line)

	if err != nil {
		return "No wallpaper fetched yet"
	}

	return filename
}

func SetSelectedWallpaper(selected string, app *appcontext.AppContext) (string, error) {
	output, err := files.ApplyWallpaper(selected, app)
	if err != nil {
		return "", err
	}

	currentID := parseIdFromPath(selected)
	wallhaven_url := fmt.Sprintf("https://wallhaven.cc/w/%s", currentID)

	//TODO: send Id to /tmp/wallchemy.sock if it exists (for wallchemy-sync) only if fromsync flag has been set

	fromSync := app.Config.GetBool("fromsync")
	if !fromSync {
		_, err := network.SendIPCMessage(currentID)
		if err != nil {

		}
	}

	current_string := fmt.Sprintf("%s\n%s\n%s\n%s", currentID, wallhaven_url, selected, output)

	err = app.CacheTools.WriteStringToFile("wallhaven/current", current_string)
	if err != nil {
		return "", err
	}
	return output, nil
}

func processQuery(app *appcontext.AppContext) error {
	last, err := processPage(app)
	if err != nil {
		return err
	}

	if last > 1 {
		lastPage := min(last, app.Config.GetIntWithDefault("max_pages", 5))
		if app.Config.GetString("collection") != "" {
			lastPage = last
		}

		var wg sync.WaitGroup
		errChan := make(chan error, lastPage-1)

		for page := 2; page <= lastPage; page++ {
			wg.Add(1)
			go func(p int) {
				defer wg.Done()

				localApp := *app
				localApp.URLBuilder = app.URLBuilder.Clone()
				localApp.URLBuilder.SetInt("page", p)

				if _, err := processPage(&localApp); err != nil {
					errChan <- fmt.Errorf("error grabbing page %d: %w", p, err)
				}
			}(page)
		}

		wg.Wait()
		close(errChan)

		if len(errChan) > 0 {
			logger.Log.Error(<-errChan)
		}
	}

	return nil
}

func processPage(app *appcontext.AppContext) (int, error) {
	request := app.URLBuilder.Build()
	logger.Log.WithField("request", request).Debug("Proccessing page")

	resp, err := download.FetchJson(request)
	if err != nil {
		return 0, fmt.Errorf("Could not fetch page: %w", err)
	}

	var wd WallhavenRequest
	if err := json.Unmarshal(resp, &wd); err != nil {
		return 0, fmt.Errorf("Could not process JSON data: %w", err)
	}

	var links []string
	for _, link := range wd.Wallpapers {
		links = append(links, link.Path)
	}

	app.LinkManager.AddLinks(links)

	return wd.Meta.LastPage, nil
}

func parseIdFromPath(path string) string {
	re := regexp.MustCompile(`wallhaven-([a-z0-9]+)\.(png|jpg|jpeg)`)
	matches := re.FindStringSubmatch(path)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}
