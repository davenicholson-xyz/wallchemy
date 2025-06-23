package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	wapi "github.com/davenicholson-xyz/go-wallhaven/wallhavenapi"
	"github.com/davenicholson-xyz/wallchemy/app"
)

type WallchemyCommand string

const (
	SingleWallpaper WallchemyCommand = "single"
	Collection      WallchemyCommand = "collection"
	CollectionList  WallchemyCommand = "collections"
	Random          WallchemyCommand = "random"
	Info            WallchemyCommand = "info"
	Clean           WallchemyCommand = "clean"
	Daemon          WallchemyCommand = "daemon"
	LaunchDaemonCmd WallchemyCommand = "launchdaemon"
	Blacklist       WallchemyCommand = "blacklist"
	ListBlacklist   WallchemyCommand = "listblacklist"
)

func globalFilters(appCtx *app.AppContext) error {

	if len(appCtx.Config.Categories) > 0 {
		var flags []wapi.CategoriesFlag
		for _, cat := range appCtx.Config.Categories {
			switch strings.ToLower(cat) {
			case "general":
				flags = append(flags, wapi.General)
			case "anime":
				flags = append(flags, wapi.Anime)
			case "people":
				flags = append(flags, wapi.People)
			default:
				return fmt.Errorf("invalid category: %s", cat)
			}
		}
		if len(flags) > 0 {
			appCtx.WallhavenAPI.Categories(flags...)
		}
	}

	if len(appCtx.Config.Purity) > 0 {
		var flags []wapi.PurityFlag
		for _, purity := range appCtx.Config.Purity {
			switch strings.ToLower(purity) {
			case "sfw":
				flags = append(flags, wapi.SFW)
			case "sketchy":
				flags = append(flags, wapi.Sketchy)
			case "nsfw":
				flags = append(flags, wapi.NSFW)
			default:
				return fmt.Errorf("invalid purity: %s", purity)
			}
		}
		if len(flags) > 0 {
			appCtx.WallhavenAPI.Purity(flags...)
		}
	}

	if appCtx.Config.TopRange != "" {
		var rangeType wapi.RangeType
		switch appCtx.Config.TopRange {
		case "1d":
			rangeType = wapi.OneDay
		case "3d":
			rangeType = wapi.ThreeDays
		case "1w":
			rangeType = wapi.OneWeek
		case "1M":
			rangeType = wapi.OneMonth
		case "3M":
			rangeType = wapi.ThreeMonths
		case "6M":
			rangeType = wapi.SixMonths
		case "1y":
			rangeType = wapi.OneYear
		default:
			return fmt.Errorf("invalid range: %s", appCtx.Config.TopRange)
		}
		appCtx.WallhavenAPI.Range(rangeType)
	}

	if appCtx.Config.MinResolution != "" {
		appCtx.WallhavenAPI.MinimumResolution(appCtx.Config.MinResolution)
	}

	if len(appCtx.Config.Resolutions) > 0 {
		appCtx.WallhavenAPI.Resolutions(appCtx.Config.Resolutions...)
	}

	if len(appCtx.Config.Ratios) > 0 {
		appCtx.WallhavenAPI.Ratios(appCtx.Config.Ratios...)
	}

	if appCtx.Config.Colors != "" {
		appCtx.WallhavenAPI.Colors(appCtx.Config.Colors)
	}

	if appCtx.Config.Seed != "" {
		appCtx.WallhavenAPI.Seed(appCtx.Config.Seed)
	}

	return nil
}

func parseAction(appCtx *app.AppContext) WallchemyCommand {

	if appCtx.Config.ID != "" {
		return SingleWallpaper
	}

	if appCtx.Config.Collection != "" {
		return Collection
	}

	if appCtx.Config.Collections {
		return CollectionList
	}

	if appCtx.Config.Random != "" || appCtx.Config.Hot || appCtx.Config.TopList {
		return Random
	}

	if appCtx.Config.Info {
		return Info
	}

	if appCtx.Config.Clean {
		return Clean
	}

	if appCtx.Config.LaunchDaemon {
		return LaunchDaemonCmd
	}

	if appCtx.Config.Daemon != "" {
		return Daemon
	}

	if appCtx.Config.Blacklist != "" {
		return Blacklist
	}

	if appCtx.Config.ListBlacklist {
		return ListBlacklist
	}

	return ""

}

func ExecuteCommand(appCtx *app.AppContext) (string, error) {

	action := parseAction(appCtx)
	app.Logger.WithField("action", string(action)).Debug("Action received")

	switch action {
	case SingleWallpaper:
		return handleSingle(appCtx)
	case Collection:
		return handleCollection(appCtx)
	case CollectionList:
		return handleCollections(appCtx)
	case Random:
		return handleRandom(appCtx)
	case Clean:
		return handleClean(appCtx)
	case Info:
		return handleInfo(appCtx)
	case LaunchDaemonCmd:
		return startDaemon(appCtx)
	case Daemon:
		return HandleDaemon(appCtx)
	case Blacklist:
		return handleBlacklist(appCtx)
	case ListBlacklist:
		return handleListBlacklist(appCtx)
	default:
		return "", fmt.Errorf("no valid action found")
	}

}

func handleClean(appCtx *app.AppContext) (string, error) {
	appCtx.CacheTools.Clean()
	app.Logger.WithField("cache dir", appCtx.CacheTools.GetCacheDirectory()).Debug("action received")
	return appCtx.CacheTools.GetCacheDirectory() + " cleaned", nil
}

func handleBlacklist(appCtx *app.AppContext) (string, error) {
	wallpaperID := appCtx.Config.Blacklist
	
	if wallpaperID == "current" {
		currentWallpaperID, err := getCurrentWallpaperID(appCtx)
		if err != nil {
			return "", fmt.Errorf("failed to get current wallpaper ID: %w", err)
		}
		wallpaperID = currentWallpaperID
	}
	
	if err := appCtx.CacheTools.AddToBlacklist(wallpaperID); err != nil {
		app.Logger.WithError(err).Error("Failed to add wallpaper to blacklist")
		return "", fmt.Errorf("failed to add wallpaper %s to blacklist: %w", wallpaperID, err)
	}
	
	app.Logger.WithField("wallpaper_id", wallpaperID).Info("Added wallpaper to blacklist")
	return fmt.Sprintf("Added wallpaper %s to blacklist", wallpaperID), nil
}

func handleListBlacklist(appCtx *app.AppContext) (string, error) {
	blacklist, err := appCtx.CacheTools.GetBlacklist()
	if err != nil {
		app.Logger.WithError(err).Error("Failed to retrieve blacklist")
		return "", fmt.Errorf("failed to retrieve blacklist: %w", err)
	}
	
	if len(blacklist) == 0 {
		return "No wallpapers in blacklist", nil
	}
	
	result := fmt.Sprintf("Blacklisted wallpapers (%d):\n", len(blacklist))
	for _, id := range blacklist {
		result += fmt.Sprintf("- %s\n", id)
	}
	
	app.Logger.WithField("blacklist_count", len(blacklist)).Info("Listed blacklisted wallpapers")
	return result, nil
}

func getCurrentWallpaperID(appCtx *app.AppContext) (string, error) {
	if !appCtx.CacheTools.FileExists("current") {
		return "", fmt.Errorf("no current wallpaper found")
	}

	infoBytes, err := appCtx.CacheTools.ReadFileBytes("current")
	if err != nil {
		return "", fmt.Errorf("failed to read current wallpaper info: %w", err)
	}

	var current app.CurrentWallpaper
	if err := json.Unmarshal(infoBytes, &current); err != nil {
		return "", fmt.Errorf("failed to parse current wallpaper info: %w", err)
	}

	if current.ID == "" {
		return "", fmt.Errorf("current wallpaper has no ID")
	}

	return current.ID, nil
}
