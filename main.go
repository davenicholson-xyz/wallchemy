package main

import (
	"fmt"
	"os"

	wapi "github.com/davenicholson-xyz/go-wallhaven/wallhavenapi"
	"github.com/davenicholson-xyz/wallchemy/app"
	"github.com/davenicholson-xyz/wallchemy/handlers"
)

var version = "version"

func main() {

	appCtx := app.NewAppContext(version)

	logConfig := app.LoggerConfigFromEnv()
	if err := app.InitLogger(logConfig); err != nil {
		panic("Failed to initialize logging")
	}

	app.Logger.Info("Wallchemy started")

	config, err := app.LoadConfig()
	if err != nil {
		app.Logger.WithError(err).Error("Unable to load config")
	}

	appCtx.AddConfig(config)

	if config.Version {
		app.Logger.WithField("version", version).Info("Version check")
		fmt.Println(version)
		os.Exit(0)
	}

	if config.Help {
		app.Logger.Info("Help requested")
		app.ShowHelp()
		os.Exit(0)
	}

	cache, err := app.NewCacheTools("wallchemy")
	if err != nil {
		panic("Failed to parse cache folder")
	}
	app.Logger.WithField("folder", cache.GetCacheDirectory()).Info("Cache folder initialized")

	if err := runApp(appCtx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func runApp(appCtx *app.AppContext) error {

	var result string

	appCtx.AddWallhavenAPI(wapi.New())

	cache, err := app.NewCacheTools("wallchemy")
	if err != nil {
		return fmt.Errorf("Cache error: %w", err)
	}
	appCtx.AddCacheTools(cache)

	if appCtx.Config.APIKey != "" {
		appCtx.WallhavenAPI.APIKey(appCtx.Config.APIKey)
	}

	result, err = handlers.ExecuteCommand(appCtx)
	if err != nil {
		return err
	}

	if appCtx.Config.Silent {
		app.Logger.WithField("output", result).Debug("Silencing output")
	} else {
		if result != "" {
			fmt.Println(result)
		}
	}

	return nil
}
