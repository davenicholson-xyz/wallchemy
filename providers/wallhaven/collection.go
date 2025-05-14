package wallhaven

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/download"
	"github.com/davenicholson-xyz/wallchemy/files"
)

type CollectionResponse struct {
	Data []Collection `json:"data"`
}

type Collection struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

func SelectionFromCollection(app *appcontext.AppContext) (string, error) {
	lm := download.NewLinkManager()
	app.AddLinkManager(lm)

	slog.Debug("Selection from collection")

	url := download.NewURL("https://wallhaven.cc/api/v1/collections")
	app.AddURLBuilder(url)

	username := app.Config.GetString("username")
	if username == "" {
		return "", fmt.Errorf("%w", errors.New("username is required for collections"))
	}

	apikey := app.Config.GetString("apikey")
	if apikey == "" {
		return "", fmt.Errorf("%w", errors.New("apikey is required for collections"))
	}

	app.URLBuilder.AddString("apikey", apikey)

	collectionID, err := getCollectionID(app)
	if err != nil {
		return "", err
	}

	if collectionID < 1 {
		return "", fmt.Errorf("%w", errors.New("Collection not found"))
	}

	err = getCollection(collectionID, app)
	if err != nil {
		return "", err
	}

	selected, err := files.GetRandomLine(app.CacheTools.Join("wallhaven", strconv.Itoa(collectionID)))
	if err != nil {
		return "", err
	}

	return selected, nil
}

func getCollectionID(app *appcontext.AppContext) (int, error) {
	slog.Debug("Getting selected collection ID")

	if !(files.PathExists(app.CacheTools.Join("wallhaven", "collections"))) {
		slog.Debug("wallhaven/collections does not exist... fetching")
		resp, err := download.FetchJson(app.URLBuilder.Build())
		if err != nil {
			return 0, err
		}

		err = app.CacheTools.WriteStringToFile("wallhaven/collections", string(resp))
		if err != nil {
			return 0, err
		}
	}

	rawJSON, err := app.CacheTools.ReadFile("wallhaven/collections")
	if err != nil {
		return 0, err
	}

	var cr CollectionResponse
	if err := json.Unmarshal([]byte(rawJSON), &cr); err != nil {
		return 0, err
	}

	collections := make(map[string]int)
	for _, coll := range cr.Data {
		collections[coll.Label] = coll.ID
	}

	collectionID := collections[app.Config.GetString("collection")]

	return collectionID, nil
}

func getCollection(collectionID int, app *appcontext.AppContext) error {

	if files.IsFileFresh(app.CacheTools.Join("wallhaven", strconv.Itoa(collectionID)), app.Config.GetIntWithDefault("expiry", 600)) {
		slog.Debug("collection file is FRESH")
		return nil
	}

	username := app.Config.GetString("username")
	apikey := app.Config.GetString("apikey")

	collectionURL := fmt.Sprintf("https://wallhaven.cc/api/v1/collections/%s/%d", username, collectionID)
	colURL := download.NewURL(collectionURL)
	colURL.AddString("apikey", apikey)
	app.AddURLBuilder(colURL)

	err := processQuery(app)
	if err != nil {
		return err
	}

	all_links := strings.Join(app.LinkManager.GetLinks(), "\n")
	err = app.CacheTools.WriteStringToFile(filepath.Join("wallhaven", strconv.Itoa(collectionID)), all_links)
	if err != nil {
		return err
	}

	return nil
}
