package chain

import (
	"context"
	"github.com/filecoin-project/lotus/chain/types"
)

type Watcher interface {
	Watch(ctx context.Context) chan types.TipSetKey
}
