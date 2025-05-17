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

var version = "version"

func main() {

	if os.Getenv("DEBUG") != "" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelError)
	}

	flg := config.NewFlagSet()

	flg.DefineString("config", "", "Path to config file")
	flg.DefineString("provider", "wallhaven", "Wallpaper provider [wallhaven]")

	flg.DefineString("username", "", "wallhaven.cc username")
	flg.DefineString("apikey", "", "wallhaven.cc api key")
	flg.DefineBool("sketchy", false, "Fetch sketchy images")
	flg.DefineBool("nsfw", false, "Fetch NSFW images")
	flg.DefineBool("noanime", false, "Dont fetch anima wallpaper")
	flg.DefineString("seed", "", "Random seed for search")
	flg.DefineInt("expiry", 0, "Cache expiry in seconds")
	flg.DefineString("extras", "", "Add the the query string")

	flg.DefineString("collection", "", "Users collection to select from")
	flg.DefineString("random", "", "Query for random wallpaper")
	flg.DefineBool("hot", false, "Select random wallpaper from the hot list")
	flg.DefineBool("top", false, "Select random wallpaper from the toplist")
	flg.DefineString("id", "", "Fetch wallpaper with ID")

	flg.DefineString("script", "", "Custom script to set wallpaper")
	flg.DefineBool("stealth", false, "Fetch but do not set wallpaper")
	flg.DefineBool("silent", false, "Do not output filename after fetch")
	flg.DefineBool("file", false, "Output the filename of the current wallpaper")
	flg.DefineBool("url", false, "Output the url of the current wallpaper")
	flg.DefineBool("link", false, "Output the provider link of the current wallpaper")
	// TODO: last - return last random query

	flg.DefineBool("clear", false, "Clear the wallchemy cache")

	flg.DefineBool("version", false, "Version info")

	flgValues := flg.Collect()

	_, version_ok := flgValues["version"]

	if version_ok {
		fmt.Println(version)
		os.Exit(0)
	}

	result, err := RunApp(flgValues)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if result != "" {
		fmt.Println(result)
	}

}

func RunApp(flgValues map[string]any) (string, error) {
	app := appcontext.NewAppContext()

	var configPath string
	cfgVal, ok := flgValues["config"].(string)
	if ok {
		configPath = cfgVal
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
