package external

import (
	"context"
	"github.com/filecoin-project/lotus/api/v0api"
	cliutil "github.com/filecoin-project/lotus/cli/util"
	"github.com/pando-project/fil-chain-extractor/pkg/lfx"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

func InjectLotusFullNode(cCtx *cli.Context) lfx.Option {
	return lfx.Override(new(v0api.FullNode), func(lc fx.Lifecycle) (v0api.FullNode, error) {
		fullNode, closeFullNode, err := cliutil.GetFullNodeAPI(cCtx)
		if err != nil {
			return nil, err
		}
		lc.Append(fx.Hook{
			OnStop: func(context.Context) error {
				closeFullNode()
				return nil
			},
		})

		return fullNode, nil
	})
}
