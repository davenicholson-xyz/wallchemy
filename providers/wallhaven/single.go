package wallhaven

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/download"
)

type SingleResponse struct {
	Data WallpaperData `json:"data,omitempty"`
}

func SelectionFromID(app *appcontext.AppContext) (string, error) {

	id := app.Config.GetString("id")

	query_url := fmt.Sprintf("https://wallhaven.cc/api/v1/w/%s", id)
	url := download.NewURL(query_url)
	app.AddURLBuilder(url)

	apikey := app.Config.GetString("apikey")
	if apikey != "" {
		app.URLBuilder.AddString("apikey", apikey)
	}

	q := url.Build()

	resp, err := download.FetchJson(q)
	if err != nil {
		if strings.Contains(err.Error(), "Error:404") {
			return "", fmt.Errorf("Wallpaper does not exist with that ID")
		}
		if strings.Contains(err.Error(), "Error:401") {
			return "", fmt.Errorf("API key required to access this wallpaper")
		}
		return "", fmt.Errorf("request failed: %w", err)
	}

	var singleResp SingleResponse
	if err := json.Unmarshal(resp, &singleResp); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	selected := singleResp.Data.Path

	return selected, nil
}

