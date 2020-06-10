package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/philips-labs/tabia/cmd"
)

const (
	appName = "tabia"
)

var (
	version = "dev"
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.HelpName = appName
	app.Usage = "code characteristics insights"
	app.EnableBashCompletion = true
	app.Version = version

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s %s/%s\n", appName, app.Version, runtime.GOOS, runtime.GOARCH)
	}

	app.Commands = cmd.CreateCommands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
