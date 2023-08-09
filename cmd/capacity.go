package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/pando-project/fil-chain-extractor/config"
	"github.com/pando-project/fil-chain-extractor/config/serialize"
	"github.com/pando-project/fil-chain-extractor/pkg/chain"
	"github.com/pando-project/fil-chain-extractor/pkg/core"
	"github.com/pando-project/fil-chain-extractor/pkg/external"
	"github.com/pando-project/fil-chain-extractor/pkg/extractor"
	"github.com/pando-project/fil-chain-extractor/pkg/lfx"
	storageMongo "github.com/pando-project/fil-chain-extractor/pkg/storage/mongo"
	"github.com/pando-project/fil-chain-extractor/pkg/util/tools"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

const (
	GENSISTIME int64 = 1598236800
	BLOCKTIME  int64 = 30
)

func NewDeltaCmd() *cli.Command {
	return &cli.Command{
		Name:  "delta",
		Usage: "",
		Subcommands: []*cli.Command{
			CapacityCmd,
			CapacityTimeCmd,
			CapacityTimeRangeCmd,
		},
	}
}

var CapacityCmd = &cli.Command{
	Name:  "capacity",
	Usage: "Latest Raw byte capacity of a minerID",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "repo-path",
			Usage: "daemon will read/persistent configs and data from the repository path",
		},
		&cli.StringFlag{
			Name:     "miner-id",
			Usage:    "Get the latest capacity of minerID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		// repo
		repoPath := c.String("repo-path")
		configPath, err := config.Filename(repoPath, "")
		if err != nil {
			logger.Errorf("failed to get config file path: %s", err)
			return err
		}
		// minerId
		minerId := c.String("miner-id")
		address, err := address.NewFromString(minerId)
		if err != nil {
			logger.Errorf("address string to go-address failed: %s", err)
			return err
		}
		// di
		ctx := context.Background()
		var modules struct {
			fx.In
			Config          *config.Config
			Storage         *storageMongo.DB
			TipSetExtractor extractor.Extractor
			Watcher         chain.Watcher
		}
		stopDeltaCapacity, err := lfx.New(ctx,
			lfx.Override(new(serialize.ConfigPath), serialize.ConfigPath(configPath)),
			external.InjectLotusFullNode(c),
			core.NewCore(ctx, core.LfxLog, &modules),
		)
		defer stopDeltaCapacity(c.Context)
		if err != nil {
			logger.Errorf("failed to inject modules: %s", err)
			return err
		}
		// power
		power, height := modules.TipSetExtractor.ExtractCapacityByMinerID(ctx, address)
		fmt.Println("miner power: ", power, "height: ", height)
		return nil
	},
}

var CapacityTimeCmd = &cli.Command{
	Name:  "capacityTime",
	Usage: "Raw byte capacity of a minerID at a given point in time",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "repo-path",
			Usage: "daemon will read/persistent configs and data from the repository path",
		},
		&cli.StringFlag{
			Name:     "unix-time",
			Usage:    "Unix timestamp",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "miner-id",
			Usage:    "Get the latest capacity of minerID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		// repo
		repoPath := c.String("repo-path")
		configPath, err := config.Filename(repoPath, "")
		if err != nil {
			logger.Errorf("failed to get config file path: %s", err)
			return err
		}
		// minerId
		minerId := c.String("miner-id")
		addr, err := address.NewFromString(minerId)
		if err != nil {
			logger.Errorf("address string to go-address failed: %s", err)
			return err
		}
		// di
		ctx := context.Background()
		var modules struct {
			fx.In
			Config          *config.Config
			Storage         *storageMongo.DB
			TipSetExtractor extractor.Extractor
			Watcher         chain.Watcher
		}
		stopDeltaCapacity, err := lfx.New(ctx,
			lfx.Override(new(serialize.ConfigPath), serialize.ConfigPath(configPath)),
			external.InjectLotusFullNode(c),
			core.NewCore(ctx, core.LfxLog, &modules),
		)
		defer stopDeltaCapacity(c.Context)
		if err != nil {
			logger.Errorf("failed to inject modules: %s", err)
			return err
		}
		// unix time
		searchTime := c.String("unix-time")
		searchTime = "1691337600"
		// genesisTime
		genesisTime := GENSISTIME
		// blockTime
		blockTime := BLOCKTIME
		// compute epoch
		epoch := tools.TimeToEpoch(searchTime, genesisTime, blockTime)
		fmt.Println("epoch: ", epoch)
		//epoch = 3103853
		// power
		power := modules.TipSetExtractor.ExtractCapacityByMinerIDAndTime(ctx, addr, abi.ChainEpoch(epoch))
		fmt.Println(power)
		return nil
	},
}

var CapacityTimeRangeCmd = &cli.Command{
	Name:  "timeRange",
	Usage: "Amount of data sealed in a time range",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "repo-path",
			Usage: "daemon will read/persistent configs and data from the repository path",
		},
		&cli.StringFlag{
			Name:     "start-utime",
			Usage:    "Start Unix timestamp",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "end-utime",
			Usage:    "End Unix timestamp",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "miner-id",
			Usage:    "Get the latest capacity of minerID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		// repo
		repoPath := c.String("repo-path")
		configPath, err := config.Filename(repoPath, "")
		if err != nil {
			logger.Errorf("failed to get config file path: %s", err)
			return err
		}
		// minerId
		minerId := c.String("miner-id")
		addr, err := address.NewFromString(minerId)
		if err != nil {
			logger.Errorf("address string to go-address failed: %s", err)
			return err
		}
		// di
		ctx := context.Background()
		var modules struct {
			fx.In
			Config          *config.Config
			Storage         *storageMongo.DB
			TipSetExtractor extractor.Extractor
			Watcher         chain.Watcher
		}
		stopDeltaCapacity, err := lfx.New(ctx,
			lfx.Override(new(serialize.ConfigPath), serialize.ConfigPath(configPath)),
			external.InjectLotusFullNode(c),
			core.NewCore(ctx, core.LfxLog, &modules),
		)
		defer stopDeltaCapacity(c.Context)
		if err != nil {
			logger.Errorf("failed to inject modules: %s", err)
			return err
		}
		// start unix time 1691337600
		start := c.String("start-utime")
		start = "1691337600"
		// end unix time 1691341200
		end := c.String("end-utime")
		end = "1691370000"

		// genesisTime
		genesisTime := GENSISTIME
		// blockTime
		blockTime := BLOCKTIME
		// compute epoch
		startEpoch := tools.TimeToEpoch(start, genesisTime, blockTime)
		endEpoch := tools.TimeToEpoch(end, genesisTime, blockTime)
		fmt.Println("start epoch: ", startEpoch, "end epoch: ", endEpoch)
		// power
		startPower := modules.TipSetExtractor.ExtractCapacityByMinerIDAndTime(ctx, addr, abi.ChainEpoch(startEpoch))
		endPower := modules.TipSetExtractor.ExtractCapacityByMinerIDAndTime(ctx, addr, abi.ChainEpoch(endEpoch))
		fmt.Println(startPower)
		fmt.Println(endPower)
		// change power
		change := big.Sub(endPower.MinerPower.RawBytePower, startPower.MinerPower.RawBytePower)
		fmt.Println("power change: ", change)
		return nil
	},
}
