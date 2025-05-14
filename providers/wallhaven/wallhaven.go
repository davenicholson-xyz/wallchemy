package wallhaven

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/download"
	"github.com/davenicholson-xyz/wallchemy/files"
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
	app.Config.Override("verbose", true)
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
		last_page := min(last, app.Config.GetIntWithDefault("max_pages", 5))
		if app.Config.GetString("collection") != "" {
			last_page = last
		}
		for page := 2; page <= last_page; page++ {
			app.URLBuilder.SetInt("page", page)
			_, err = processPage(app)
			if err != nil {
			}
		}
	}
	return nil
}

func processPage(app *appcontext.AppContext) (int, error) {
	request := app.URLBuilder.Build()

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
