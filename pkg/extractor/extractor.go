package extractor

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/types"
	minerPower "github.com/pando-project/fil-chain-extractor/pkg/schema/MinerPower"
)

type Extractor interface {
	Extract(context.Context, chan types.TipSetKey) chan any
	ExtractThenPersist(context.Context, chan types.TipSetKey)
	ExtractCapacityByMinerID(context.Context, address.Address) (*minerPower.MinerPower, abi.ChainEpoch)
	ExtractCapacityByMinerIDAndTime(context.Context, address.Address, abi.ChainEpoch) *minerPower.MinerPower
}
