package app

import (
	wapi "github.com/davenicholson-xyz/go-wallhaven/wallhavenapi"
)

type CurrentWallpaper struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Link     string `json:"link"`
	FilePath string `json:"filepath"`
}

type AppContext struct {
	Config       *Config
	CacheTools   *CacheTools
	WallhavenAPI *wapi.WallhavenAPI
	Version      string
	Current      CurrentWallpaper
}

func NewAppContext(version string) *AppContext {
	return &AppContext{Version: version, Current: CurrentWallpaper{}}
}

func (app *AppContext) AddConfig(config *Config) {
	app.Config = config
}

func (app *AppContext) AddWallhavenAPI(wapi *wapi.WallhavenAPI) {
	app.WallhavenAPI = wapi
}

func (app *AppContext) AddCacheTools(cachetools *CacheTools) {
	app.CacheTools = cachetools
}
