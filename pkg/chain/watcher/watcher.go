package watcher

import (
	"context"
	"github.com/filecoin-project/lotus/chain/store"
	"time"

	"github.com/filecoin-project/lotus/api/v0api"
	"github.com/filecoin-project/lotus/chain/types"

	"github.com/pando-project/fil-chain-extractor/config"
	"github.com/pando-project/fil-chain-extractor/pkg/util/log"
)

var logger = log.NewSubsystemLogger()

type Config struct {
	minWatchInterval time.Duration
	maxWatchInterval time.Duration
}

type Watcher struct {
	LotusFullNode v0api.FullNode
	Interval      time.Duration
	Config        Config
}

func NewWatcher(lotusFullNode v0api.FullNode, conf *config.Config) *Watcher {
	return &Watcher{
		LotusFullNode: lotusFullNode,
		Config: Config{
			minWatchInterval: conf.Watcher.MinWatchInterval,
			maxWatchInterval: conf.Watcher.MaxWatchInterval,
		},
	}
}

func (w *Watcher) Watch(ctx context.Context) chan types.TipSetKey {
	tipSetKey := make(chan types.TipSetKey, 1)
	go w.watch(ctx, tipSetKey)
	return tipSetKey
}

func (w *Watcher) watch(ctx context.Context, tipSetKey chan types.TipSetKey) {
	logger.Infof("start watching the latest chain header...")
	defer logger.Infof("chain header watch terminated.")

	cancel := context.CancelFunc(func() {})
	chainHeadChangeCh, err := w.LotusFullNode.ChainNotify(ctx)
	if err != nil {
		logger.Fatalf("failed to get chain notify channel: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			cancel()
			return
		case chainHeadChange, ok := <-chainHeadChangeCh:
			if !ok {
				logger.Errorf("failed to get chain head update")
				cancel()
				return
			}
			for _, chainHeadUpdate := range chainHeadChange {
				logger.Infof("=========================================")
				logger.Infof("get chain header:")
				logger.Infof("\ttype: %s", chainHeadUpdate.Type)
				logger.Infof("\tkey: %s", chainHeadUpdate.Val.Key().String())
				switch chainHeadUpdate.Type {
				case store.HCApply, store.HCCurrent:
					tipSetKey <- chainHeadUpdate.Val.Key()
				}
				logger.Infof("\tcids: %s", chainHeadUpdate.Val.Cids())
				logger.Infof("\theight: %s", chainHeadUpdate.Val.Height().String())
			}
		}
	}
}
