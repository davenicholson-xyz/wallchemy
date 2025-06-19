package handlers

import (
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
	Current         WallchemyCommand = "current"
	Clean           WallchemyCommand = "clean"
	Daemon          WallchemyCommand = "daemon"
	LaunchDaemonCmd WallchemyCommand = "launchdaemon"
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

	if appCtx.Config.Current {
		return Current
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
	case Current:
		return handleCurrent(appCtx)
	case LaunchDaemonCmd:
		return startDaemon(appCtx)
	case Daemon:
		return HandleDaemon(appCtx)
	default:
		return "", fmt.Errorf("no valid action found")
	}

}

func handleClean(appCtx *app.AppContext) (string, error) {
	appCtx.CacheTools.Clean()
	app.Logger.WithField("cache dir", appCtx.CacheTools.GetCacheDirectory()).Debug("action received")
	return appCtx.CacheTools.GetCacheDirectory() + " cleaned", nil
}
