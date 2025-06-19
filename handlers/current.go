package handlers

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/davenicholson-xyz/wallchemy/app"
)

func handleCurrent(appCtx *app.AppContext) (string, error) {

	if appCtx.CacheTools.FileExists("current") {

		infoStr, err := appCtx.CacheTools.ReadFileBytes("current")
		if err != nil {
			return "", err
		}

		var current app.CurrentWallpaper
		if err := json.Unmarshal(infoStr, &current); err != nil {
			return "", err
		}

		return fmt.Sprintf("%s\n%s\n%s\n%s", current.ID, current.Path, current.Link, current.FilePath), nil

	}

	return "", fmt.Errorf("No current wallpaper info")

}

func SaveCurrent(appCtx *app.AppContext) error {
	fileName := filepath.Base(appCtx.Current.Path)

	id, err := parseIdFromFileName(fileName)
	if err != nil {
		return err
	}

	appCtx.Current.ID = id
	appCtx.Current.Link = fmt.Sprintf("https://wallhaven.cc/w/%s", id)

	current, err := json.MarshalIndent(appCtx.Current, "", " ")
	if err != nil {
		return err
	}

	if err = appCtx.CacheTools.WriteStringToFile("current", string(current)); err != nil {
		return err
	}

	app.Logger.Debug("Saved current file info")

	return nil
}

func parseIdFromFileName(fileName string) (string, error) {

	if len(fileName) > 10 && fileName[:10] == "wallhaven-" {
		fileName = fileName[10:]

		if ext := filepath.Ext(fileName); ext != "" {
			fileName = fileName[:len(fileName)-len(ext)]
			return fileName, nil
		}
	}

	return "", fmt.Errorf("could not extract ID")

}
