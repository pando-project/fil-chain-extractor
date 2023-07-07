package main

import (
	"context"
	"github.com/pando-project/fil-chain-extractor/pkg/extractor"
	storageMongo "github.com/pando-project/fil-chain-extractor/pkg/storage/mongo"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"

	"github.com/pando-project/fil-chain-extractor/config"
	"github.com/pando-project/fil-chain-extractor/config/serialize"
	"github.com/pando-project/fil-chain-extractor/pkg/chain"
	"github.com/pando-project/fil-chain-extractor/pkg/core"
	"github.com/pando-project/fil-chain-extractor/pkg/external"
	"github.com/pando-project/fil-chain-extractor/pkg/lfx"
)

func NewDaemonCmd() *cli.Command {
	return &cli.Command{
		Name:  "daemon",
		Usage: "Start a github.com/pando-project/fil-chain-extractor daemon process",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "lotus-api-token",
				Usage:    "lotus api auth token which has read permission",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "lotus-api-address",
				Usage:    "lotus full-node api address in multi-address manner",
				Value:    "/ip4/127.0.0.1/tcp/1234",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "repo-path",
				Usage:    "daemon will read/persistent configs and data from the repository path",
				Required: false,
			},
		},
		Action: func(cCtx *cli.Context) error {
			repoPath := cCtx.String("repo-path")
			configPath, err := config.Filename(repoPath, "")
			if err != nil {
				logger.Errorf("failed to get config file path: %s", err)
				return err
			}
			ctx := context.Background()
			var modules struct {
				fx.In
				Config          *config.Config
				Storage         *storageMongo.DB
				Watcher         chain.Watcher
				TipSetExtractor extractor.Extractor
			}
			stopDaemon, err := lfx.New(ctx,
				lfx.Override(new(serialize.ConfigPath), serialize.ConfigPath(configPath)),
				core.NewCore(ctx, core.LfxLog, &modules),
				external.InjectLotusFullNode(cCtx),
			)
			defer stopDaemon(cCtx.Context)
			if err != nil {
				logger.Errorf("failed to inject modules: %s", err)
				return err
			}

			tipSetKey := modules.Watcher.Watch(ctx)
			modules.TipSetExtractor.ExtractThenPersist(ctx, tipSetKey)

			return nil
		},
	}
}
