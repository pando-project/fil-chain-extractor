package main

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"

	"github.com/pando-project/fil-chain-extractor/pkg/fcelog"
)

var logger = logging.Logger("main")

func main() {
	fcelog.SetupLogLevels()

	local := []*cli.Command{
		NewDaemonCmd(),
		NewInitCmd(),
		NewDeltaCmd(),
	}

	app := &cli.App{
		Name:                 "fce",
		Usage:                "Filecoin chain data extractor",
		EnableBashCompletion: true,
		Flags:                []cli.Flag{},
		Commands:             local,
	}

	app.Setup()
	RunApp(app)
}
