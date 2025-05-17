package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/davenicholson-xyz/wallchemy/application"
	"github.com/davenicholson-xyz/wallchemy/config"
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

	result, err := application.RunApp(flgValues)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if result != "" {
		fmt.Println(result)
	}

}
