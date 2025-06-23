package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var configFile string

type Config struct {
	Version       bool     `mapstructure:"version"`
	Help          bool     `mapstructure:"help"`
	ConfigFile    string   `mapstructure:"config"`
	APIKey        string   `mapstructure:"apikey"`
	Username      string   `mapstructure:"username"`
	ID            string   `mapstructure:"id"`
	Collection    string   `mapstructure:"collection"`
	Collections   bool     `mapstructure:"collections"`
	Random        string   `mapstructure:"random"`
	Hot           bool     `mapstructure:"hot"`
	TopList       bool     `mapstructure:"top"`
	TopRange      string   `mapstructure:"range"`
	Categories    []string `mapstructure:"categories"`
	Purity        []string `mapstructure:"purity"`
	MinResolution string   `mapstructure:"minresolution"`
	Resolutions   []string `mapstructure:"resolutions"`
	Ratios        []string `mapstructure:"ratios"`
	Colors        string   `mapstructure:"colors"`
	Seed          string   `mapstructure:"seed"`
	Expiry        int      `mapstructure:"expiry"`
	MaxPages      int      `mapstructure:"pages"`
	Output        string   `mapstructure:"output"`
	Script        string   `mapstructure:"script"`
	Quiet         bool     `mapstructure:"quiet"`
	Stealth       bool     `mapstructure:"stealth"`
	Clean         bool     `mapstructure:"clean"`
	Info          bool     `mapstructure:"info"`
	Daemon        string   `mapstructure:"daemon"`
	LaunchDaemon  bool     `mapstructure:"launchdaemon"`
	Port          int      `mapstructure:"port"`
	Blacklist     string   `mapstructure:"blacklist"`
	ListBlacklist bool     `mapstructure:"listblacklist"`
}

func setupFlags() {

	pflag.BoolP("version", "v", false, "Show version information and exit")
	pflag.BoolP("help", "h", false, "Show this help message and exit")

	pflag.StringVarP(&configFile, "config", "c", "", "Path to custom configuration file (YAML)")
	pflag.StringP("apikey", "k", "", "Wallhaven.cc API key for accessing private collections and NSFW content")
	pflag.StringP("username", "u", "", "Wallhaven.cc username (required for accessing your collections)")

	pflag.StringP("id", "i", "", "Download and set a specific wallpaper by its Wallhaven ID (e.g., 'ox1234')")
	pflag.StringP("collection", "C", "", "Set a random wallpaper from one of your named collections")
	pflag.Bool("collections", false, "List all available collections for your account")
	pflag.StringP("random", "r", "", "Search term for random wallpapers (e.g., 'nature', 'cars', 'abstract')")
	pflag.BoolP("hot", "H", false, "Select from currently hot/trending wallpapers")
	pflag.BoolP("top", "t", false, "Select from top-rated wallpapers (use --range to specify time period)")

	pflag.String("range", "", "Time range for top wallpapers: 1d, 3d, 1w, 1M, 3M, 6M, 1y")
	pflag.StringSlice("categories", []string{}, "Wallpaper categories to include: general, anime, people")
	pflag.StringSlice("purity", []string{}, "Content purity levels: sfw, sketchy, nsfw (requires API key for nsfw)")

	pflag.String("minresolution", "", "Minimum resolution filter (e.g., '1920x1080', '2560x1440')")
	pflag.StringSlice("resolutions", []string{}, "Exact resolution matches (e.g., '1920x1080,2560x1440,3840x2160')")
	pflag.StringSlice("ratios", []string{}, "Aspect ratio filters (e.g., '16x9,16x10,21x9,4x3,5x4')")

	pflag.String("colors", "", "Dominant color filter as hex code without # (e.g., 'ff0000' for red, '00ff00' for green)")
	pflag.String("seed", "", "Random seed for consistent results across searches (any string)")

	pflag.IntP("expiry", "e", 600, "Cache expiry time in seconds (how long to keep search results cached)")
	pflag.Int("pages", 5, "Maximum number of pages to fetch from API (24 wallpapers per page)")

	pflag.StringP("script", "s", "", "Custom script to set wallpaper instead of built-in method (receives file path as argument)")
	pflag.StringP("output", "o", "", "Output folder for wallpaper downloads")
	pflag.BoolP("quiet", "q", false, "Suppress all output except errors")
	pflag.Bool("stealth", false, "Download wallpaper but don't actually set it (useful for testing)")

	pflag.BoolP("clean", "x", false, "Clean the cache directory (removes all cached wallpapers and search results)")
	pflag.BoolP("info", "I", false, "Show information about the currently set wallpaper")

	pflag.StringP("daemon", "d", "", "Start the background daemon for browser extension support")
	pflag.Bool("launchdaemon", false, "")
	pflag.IntP("port", "p", 2388, "Port for the backgorund daemon to listen on.")

	pflag.StringP("blacklist", "b", "", "Add wallpaper ID to blacklist (use 'current' to blacklist currently set wallpaper)")
	pflag.Bool("listblacklist", false, "List all blacklisted wallpaper IDs")

	pflag.CommandLine.MarkHidden("launchdaemon")

	pflag.Parse()
}

func LoadConfig() (*Config, error) {

	setupFlags()

	v := viper.New()

	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("$HOME/.config/wallchemy")
		v.AddConfigPath(filepath.Join(os.Getenv("AppData"), "wallchemy"))
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Logger.Debug("Config file not found. using flags/env")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		Logger.WithField("file", v.ConfigFileUsed()).Info("Config file found")
	}

	v.SetEnvPrefix("WALLCHEMY")
	v.AutomaticEnv()

	v.BindPFlags(pflag.CommandLine)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil

}
func ShowHelp() {
	fmt.Println(`Wallchemy - Wallhaven.cc Wallpaper Manager

USAGE:
  wallchemy [OPTIONS]

WALLPAPER SELECTION MODES:
  Use one of these modes to select wallpapers:
  -i, --id <wallpaper-id>           Download specific wallpaper by ID
  -r, --random <search-term>        Search for random wallpapers
  -H, --hot                         Get hot/trending wallpapers  
  -t, --top                         Get top-rated wallpapers
  -C, --collection <name>           Random wallpaper from your collection
      --collections                 List your available collections

FILTERING OPTIONS:
      --range <period>              Time range for --top: 1d,3d,1w,1M,3M,6M,1y
      --categories <list>           Categories: general,anime,people
      --purity <list>               Content: sfw,sketchy,nsfw (nsfw needs API key)
      --minresolution <res>         Minimum resolution (e.g., 1920x1080)
      --resolutions <list>          Exact resolutions (e.g., 1920x1080,2560x1440)
      --ratios <list>               Aspect ratios (e.g., 16x9,16x10,21x9)
      --colors <hex>                Color filter (hex without #, e.g., ff0000)
      --seed <string>               Random seed for consistent results

CONFIGURATION:
  -c, --config <file>               Custom config file path
  -k, --apikey <key>                Wallhaven.cc API key
  -u, --username <name>             Wallhaven.cc username

EXECUTION:
  -s, --script <path>               Custom wallpaper-setting script
  -o, --output <path>               Output folder for wallpaper downloads
  -q, --quiet                       Suppress output
      --stealth                     Download but don't set wallpaper
  -e, --expiry <seconds>            Cache expiry (default: 600)
      --pages <num>                 Max pages to fetch (default: 5)

UTILITIES:
  -I, --info                        Show current wallpaper info
  -x, --clean                       Clean cache directory
  -b, --blacklist <id|current>      Add wallpaper ID to blacklist (use 'current' for current wallpaper)
      --listblacklist               List all blacklisted wallpapers
  -v, --version                     Show version
  -h, --help                        Show this help

BROWSER EXTENSION DAEMON:
  -d, --daemon start                Start the background daemon
  -d, --daemon status               Check background daemon is running
  -d, --daemon stop                 Stop the background daemon
  -p, --port <num>                  Port for daemon to listen on (default: 2388)

EXAMPLES:
  wallchemy --random "nature mountains"
  wallchemy --hot --categories general --purity sfw
  wallchemy --top --range 1w --minresolution 2560x1440
  wallchemy --id ox83dp3
  wallchemy --collection "my Favorites"
  wallchemy --random "cyberpunk" --colors 00ffff --ratios 21x9
  wallchemy --blacklist ox83dp3
  wallchemy --blacklist current

EXAMPLE CONFIG FILE - config.yml:
    apikey: "your-api-key-here"
    username: "yourusername"
    categories: ["general", "anime"]
    purity: ["sfw"]
    expiry: 1200
    pages: 10`)
}
