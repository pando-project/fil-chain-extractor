package config

import "time"

type Watcher struct {
	MinWatchInterval time.Duration `yaml:"MinWatchInterval"`
	MaxWatchInterval time.Duration `yaml:"MaxWatchInterval"`
}

func DefaultWatcher() Watcher {
	return Watcher{
		MinWatchInterval: time.Second,
		MaxWatchInterval: 10 * time.Second,
	}
}
