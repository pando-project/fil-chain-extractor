package main

import (
	"os"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"

	"fil-chain-extractor/pkg/fcelog"
)

var logger = logging.Logger("main")

func main() {
	fcelog.SetupLogLevels()

	local := []*cli.Command{
		NewDaemonCmd(),
		NewInitCmd(),
	}

	app := &cli.App{
		Name:                 "fce",
		Usage:                "Filecoin chain data extractor",
		EnableBashCompletion: true,
		Flags:                []cli.Flag{},
		Commands:             local,
	}

	app.Setup()
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
