package config

type LotusInfo struct {
	APIAddress string `yaml:"APIAddress"`
	Token      string `yaml:"Token"`
}

func DefaultLotusInfo() LotusInfo {
	return LotusInfo{
		APIAddress: "/ip4/127.0.0.1/tcp/1234",
	}
}
