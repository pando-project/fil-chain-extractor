package config

import (
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	API       API
	LotusInfo LotusInfo
	Storage   Storage
	Watcher   Watcher
}

const (
	DefaultRepoName   = ".fil-chain-extractor"
	DefaultRepoRoot   = "~/" + DefaultRepoName
	DefaultConfigFile = "config.yaml"
	EnvDir            = "FIL_CHAIN_EXTRACTOR_PATH"
)

func PathRoot() (string, error) {
	dir := os.Getenv(EnvDir)
	var err error
	if len(dir) == 0 {
		dir, err = homedir.Expand(DefaultRepoRoot)
	}

	return dir, err
}

func Path(configRoot, extension string) (string, error) {
	if len(configRoot) == 0 {
		dir, err := PathRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, extension), nil
	}
	return filepath.Join(configRoot, extension), nil
}

func Filename(configRoot string, userConfigFile string) (string, error) {
	if userConfigFile == "" {
		return Path(configRoot, DefaultConfigFile)
	}
	if filepath.Dir(userConfigFile) == "." {
		return Path(configRoot, userConfigFile)
	}
	return userConfigFile, nil
}

func HumanOutput(value any) ([]byte, error) {
	s, ok := value.(string)
	if ok {
		return []byte(strings.Trim(s, "\n")), nil
	}
	return Marshal(value)
}

func Marshal(value any) ([]byte, error) {
	return yaml.Marshal(value)
}
