package config

func Init() *Config {
	return &Config{
		LotusInfo: DefaultLotusInfo(),
		API:       DefaultAPI(),
		Storage:   DefaultStorage(),
		Watcher:   DefaultWatcher(),
	}
}
