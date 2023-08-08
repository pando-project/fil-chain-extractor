package tools

import (
	logging "github.com/ipfs/go-log/v2"
	"strconv"
)

var logger = logging.Logger("tools")

func TimeToEpoch(searchTime string, genesisTime, blockTime int64) int64 {
	time, err := strconv.ParseInt(searchTime, 10, 64)
	if err != nil {
		logger.Errorf("strconv string to unix-time err: %s", err)
		return -1
	}
	timeDiff := time - genesisTime
	epoch := timeDiff / blockTime
	return epoch
}
