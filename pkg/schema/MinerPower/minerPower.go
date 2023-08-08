package minerPower

import "github.com/filecoin-project/lotus/chain/actors/builtin/power"

type MinerPower struct {
	MinerPower  power.Claim
	TotalPower  power.Claim
	HasMinPower bool
}
