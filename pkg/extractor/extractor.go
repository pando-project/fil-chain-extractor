package extractor

import (
	"context"
	"github.com/filecoin-project/lotus/chain/types"
)

type Extractor interface {
	Extract(context.Context, chan types.TipSetKey) chan any
	ExtractThenPersist(context.Context, chan types.TipSetKey)
}
