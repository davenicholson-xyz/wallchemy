package application

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
