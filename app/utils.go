package app

import (
	"fmt"
	wapi "github.com/davenicholson-xyz/go-wallhaven/wallhavenapi"
	"net/url"
)

func RemoveURLParams(urlStr string, paramsToRemove ...string) (string, error) {

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	query := parsedURL.Query()

	for _, param := range paramsToRemove {
		query.Del(param)
	}

	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

func SortMethod(appCtx *AppContext) wapi.SortingType {
	if appCtx.Config.Hot {
		return wapi.Hot
	}

	if appCtx.Config.TopList {
		return wapi.Toplist
	}

	return wapi.Random
}
