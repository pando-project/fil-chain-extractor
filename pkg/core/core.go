package core

import (
	"context"
	logging "github.com/ipfs/go-log/v2"
	"github.com/pando-project/fil-chain-extractor/config"
	"github.com/pando-project/fil-chain-extractor/config/serialize"
	"github.com/pando-project/fil-chain-extractor/pkg/chain"
	"github.com/pando-project/fil-chain-extractor/pkg/chain/watcher"
	"github.com/pando-project/fil-chain-extractor/pkg/extractor"
	"github.com/pando-project/fil-chain-extractor/pkg/lfx"
	storageMongo "github.com/pando-project/fil-chain-extractor/pkg/storage/mongo"
	"github.com/pando-project/fil-chain-extractor/pkg/util/log"
	"go.uber.org/fx"
)

const (
	invokeNone lfx.Invoke = iota
	invokeSetupDebug
	invokePopulate
)

var (
	LfxLog = &LfxLogger{
		log.NewSubsystemLogger(),
	}
)

type LfxLogger struct {
	*logging.ZapEventLogger
}

func (l *LfxLogger) Printf(msg string, args ...any) {
	l.ZapEventLogger.Debugf(msg, args)
}

func NewCore(ctx context.Context, logger fx.Printer, target ...any) lfx.Option {
	return lfx.Options(
		lfx.Override(new(GlobalContext), ctx),
		lfx.If(logger != nil, lfx.Logger(logger)),
		lfx.If(len(target) > 0, lfx.Populate(invokePopulate, target...)),
		lfx.Override(new(*config.Config), serialize.Load),
		lfx.Override(new(*storageMongo.DB), storageMongo.NewMongoDB),
		lfx.Override(new(chain.Watcher), watcher.NewWatcher),
		lfx.Override(new(extractor.Extractor), extractor.NewTipSetExtractor),
	)
}
