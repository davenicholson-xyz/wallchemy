package wallhaven

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/download"
	"github.com/davenicholson-xyz/wallchemy/files"
	"github.com/davenicholson-xyz/wallchemy/logger"
)

func SelectionFromRandom(app *appcontext.AppContext) (string, error) {
	lm := download.NewLinkManager()
	app.AddLinkManager(lm)

	logger.Log.Info("Selecting from random")

	url := download.NewURL("https://wallhaven.cc/api/v1/search")
	app.AddURLBuilder(url)

	seed := app.Config.GetStringWithDefault("seed", download.GenerateSeed(6))
	app.URLBuilder.AddString("seed", seed)

	apikey := app.Config.GetString("apikey")
	if apikey != "" {
		app.URLBuilder.AddString("apikey", apikey)
	}

	url.SetString("purity", "100")
	if app.Config.GetBool("sketchy") {
		app.URLBuilder.SetString("purity", "110")
	}
	if app.Config.GetBool("nsfw") {
		app.URLBuilder.SetString("purity", "111")
	}

	url.SetString("categories", "111")
	if app.Config.GetBool("noanime") {
		app.URLBuilder.SetString("categories", "101")
	}

	random := app.Config.GetString("random")
	if random != "" {
		app.URLBuilder.SetString("sorting", "random")
		app.URLBuilder.AddString("q", random)
	}

	if app.Config.GetBool("hot") {
		url.SetString("sorting", "hot")
	}

	if app.Config.GetBool("top") {
		app.URLBuilder.SetString("sorting", "toplist")
	}

	extras := app.Config.GetString("extras")
	app.URLBuilder.AddExtras(extras)

	err := getRandom(app)
	if err != nil {
		return "", err
	}

	outfile := app.URLBuilder.GetString("sorting")
	selected, _ := files.GetRandomLine(app.CacheTools.Join("wallhaven", outfile))
	// if err != nil {
	// 	return "", err
	// }

	return selected, nil
}

func getRandom(app *appcontext.AppContext) error {

	sorting := app.URLBuilder.GetString("sorting")
	cleanUrl := app.URLBuilder.Without("apikey").Without("seed")
	query_url := cleanUrl.Build()

	slog.Debug("query_url: " + query_url)

	last_file := fmt.Sprintf("%s_query", sorting)
	if files.PathExists(app.CacheTools.Join("wallhaven", last_file)) {
		slog.Debug(last_file + " does exist")
		last_query, err := app.CacheTools.ReadLineFromFile("wallhaven/"+last_file, 1)
		slog.Debug("last_query: " + last_query)
		if err != nil {
			return err
		}
		if last_query == query_url {
			slog.Debug("last query = query_url")
			if files.IsFileFresh(app.CacheTools.Join(filepath.Join("wallhaven", sorting)), app.Config.GetIntWithDefault("expiry", 600)) {
				slog.Debug("Using cached results")
				return nil
			}
		}
	}

	lm := download.NewLinkManager()
	app.AddLinkManager(lm)

	err := processQuery(app)
	if err != nil {
		return err
	}

	if app.LinkManager.Count() == 0 {
		return errors.New("No wallpapers found")
	}

	err = app.CacheTools.WriteStringToFile(filepath.Join("wallhaven", sorting+"_query"), query_url)
	if err != nil {
		return err
	}

	all_links := app.LinkManager.GetLinks()
	slog.Debug("Links found: " + strconv.Itoa(len(all_links)))

	links_str := strings.Join(app.LinkManager.GetLinks(), "\n")
	err = app.CacheTools.WriteStringToFile(filepath.Join("wallhaven", sorting), links_str)
	if err != nil {
		return err
	}

	return nil
}

func LastRandomQuery(app *appcontext.AppContext) string {
	lastUrl, err := app.CacheTools.ReadLineFromFile("wallhaven/random_query", 1)
	if err != nil {
		logger.Log.Warn("Could not read last query for random")
	}
	parsedUrl, err := url.Parse(lastUrl)
	if err != nil {
		logger.Log.Warn("Could not parse last query for random")
	}
	queryParams := parsedUrl.Query()
	qValue := queryParams.Get("q")
	return qValue
}
