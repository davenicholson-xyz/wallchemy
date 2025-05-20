package application

import (
	"fmt"
	"os"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/cachetools"
	"github.com/davenicholson-xyz/wallchemy/config"
	"github.com/davenicholson-xyz/wallchemy/files"
	"github.com/davenicholson-xyz/wallchemy/logger"
	"github.com/davenicholson-xyz/wallchemy/providers"
)

func RunApp(flgValues map[string]any) (string, error) {
	app := appcontext.NewAppContext()

	var configPath string
	cfgVal, ok := flgValues["config"].(string)
	if ok {
		configPath = cfgVal
		if !files.PathExists(configPath) {
			logger.Log.WithField("configPath", configPath).Error("Path does not exist")
			return "", fmt.Errorf("config file does not exist")
		}
	} else {
		configPath, _ = files.DefaultConfigFilepath()
	}

	logger.Log.WithField("configPath", configPath).Debug("set config path")

	cfg, err := config.New(configPath)
	if err != nil {
		logger.Log.WithField("configPath", configPath).Error("Failed to load config file")
		return "", fmt.Errorf("Failed to load config: %w", err)
	}

	cfg.FlagOverride(flgValues)
	app.AddConfig(cfg)

	// slog.Debug(fmt.Sprintf("%+v", app.Config))

	ct, err := cachetools.New("wallchemy")
	if err != nil {
		logger.Log.Error("Failed to create cache")
		return "", fmt.Errorf("Error creating cache: %w", err)
	}

	app.AddCacheTools(ct)

	if app.Config.GetBool("clear") {
		logger.Log.Info("Clearing cache")
		err := app.CacheTools.Clear()
		if err != nil {
			logger.Log.Error("Failed to clear cache")
			return "", fmt.Errorf("Error deleting cache: %w", err)
		}
		return "Cache deleted", nil
	}

	prov := app.Config.GetString("provider")
	provider, exists := providers.GetProvider(prov)
	if !exists {
		logger.Log.WithField("provider", provider).Error("Provider not found")
		return "", fmt.Errorf("Provider not found: %q", prov)
	}

	logger.Log.Info("Checking existance of cache folder")
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
