package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/cachetools"
	"github.com/davenicholson-xyz/wallchemy/config"
	"github.com/davenicholson-xyz/wallchemy/files"
	"github.com/davenicholson-xyz/wallchemy/providers"
)

func main() {

	if os.Getenv("DEBUG") != "" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelError)
	}

	result, err := runApp()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if result != "" {
		fmt.Println(result)
	}

}

func runApp() (string, error) {
	app := appcontext.NewAppContext()

	flg := config.NewFlagSet()

	flg.DefineString("config", "", "path to config file")
	flg.DefineString("provider", "wallhaven", "wallpaper provider")

	flg.DefineString("username", "", "wallhaven.cc username")
	flg.DefineString("apikey", "", "wallhaven.cc api key")
	flg.DefineBool("nsfw", false, "Fetch NSFW images")
	flg.DefineBool("noanime", false, "dont fetch anima wallpaper")
	flg.DefineString("seed", "", "random seed for search")
	flg.DefineInt("expiry", 0, "cache expiry in seconds")

	flg.DefineString("collection", "", "collection to select from")
	flg.DefineString("random", "", "query for random wallpaper")
	flg.DefineBool("hot", false, "hot")
	flg.DefineBool("top", false, "toplist")

	flg.DefineString("script", "", "script to set wallpaper")
	flg.DefineBool("stealth", false, "do not change wallpaper")
	flg.DefineBool("silent", false, "does not output the filename")
	flg.DefineBool("file", false, "output the filename of the current wallpaper")
	flg.DefineBool("url", false, "output the url of the current wallpaper")
	flg.DefineBool("link", false, "output the provider link of the current wallpaper")

	flg.DefineBool("clear", false, "clear the wallchemy cache")

	flgValues := flg.Collect()

	var configPath string
	val, ok := flgValues["config"].(string)
	if ok {
		configPath = val
		if !files.PathExists(configPath) {
			return "", fmt.Errorf("config file does not exist")
		}
	} else {
		configPath, _ = files.DefaultConfigFilepath()
	}

	cfg, err := config.New(configPath)
	cfg.FlagOverride(flgValues)

	app.AddConfig(cfg)

	slog.Debug(fmt.Sprintf("%+v", app.Config))

	if err != nil {
		return "", fmt.Errorf("Failed to load config: %w", err)
	}

	ct, err := cachetools.New("wallchemy")
	if err != nil {
		return "", fmt.Errorf("Error creating cache: %w", err)
	}

	app.AddCacheTools(ct)

	if app.Config.GetBool("clear") {
		slog.Debug("Clearing the cache")
		err := app.CacheTools.Clear()
		if err != nil {
			return "", fmt.Errorf("Error deleting cache: %w", err)
		}
		return "Cache deleted", nil
	}

	prov := app.Config.GetString("provider")
	provider, exists := providers.GetProvider(prov)
	if !exists {
		return "", fmt.Errorf("Provider not found: %q", prov)
	}

	slog.Debug("Checking provider cache folder exists")
	err = os.MkdirAll(app.CacheTools.Join(provider.Name()), 0775)
	if err != nil {
		return "", fmt.Errorf("Error creating provider cache dir: %w", err)
	}

	result, err := provider.ParseArgs(app)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if app.Config.GetBool("silent") {
		return "", nil
	} else {
		return result, nil
	}
}
