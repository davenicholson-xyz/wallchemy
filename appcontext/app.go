package appcontext

import (
	"github.com/davenicholson-xyz/wallchemy/cachetools"
	"github.com/davenicholson-xyz/wallchemy/config"
	"github.com/davenicholson-xyz/wallchemy/download"
)

type AppContext struct {
	Config      *config.Config
	CacheTools  *cachetools.CacheTools
	URLBuilder  *download.URLBuilder
	LinkManager *download.LinkManager
}

func NewAppContext() *AppContext {
	return &AppContext{}
}

func (app *AppContext) AddConfig(cfg *config.Config) {
	app.Config = cfg
}

func (app *AppContext) AddCacheTools(ct *cachetools.CacheTools) {
	app.CacheTools = ct
}

func (app *AppContext) AddURLBuilder(url *download.URLBuilder) {
	app.URLBuilder = url
}

func (app *AppContext) AddLinkManager(lm *download.LinkManager) {
	app.LinkManager = lm
}
