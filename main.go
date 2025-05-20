package main

import (
	"fmt"
	"os"

	"github.com/davenicholson-xyz/wallchemy/application"
	"github.com/davenicholson-xyz/wallchemy/config"
	"github.com/davenicholson-xyz/wallchemy/logger"
)

var version = "version"

func main() {

	logger.Log.Info("Wallchemy starting...")

	flgValues := config.ParseFlags()

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
