package providers

import "github.com/davenicholson-xyz/wallchemy/appcontext"

type Provider interface {
	Name() string
	ParseArgs(app *appcontext.AppContext) (string, error)
}
