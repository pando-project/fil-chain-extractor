package serialize

import (
	"errors"
	"fmt"
	"github.com/facebookgo/atomicfile"
	"github.com/pando-project/fil-chain-extractor/config"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
)

// ErrNotInitialized is returned when we fail to read the config because the
// repo doesn't exist.
var ErrNotInitialized = errors.New("fce not initialized, please run 'fce init'")

type ConfigPath string

// ReadConfigFile reads the config from `filename` into `cfg`.
func ReadConfigFile(filename ConfigPath, cfg any) error {
	f, err := os.Open(string(filename))
	if err != nil {
		if os.IsNotExist(err) {
			err = ErrNotInitialized
		}
		return err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return fmt.Errorf("failure to decode config: %s", err)
	}
	return nil
}

// WriteConfigFile writes the config from `cfg` into `filename`.
func WriteConfigFile(filename ConfigPath, cfg any) error {
	err := os.MkdirAll(filepath.Dir(string(filename)), 0755)
	if err != nil {
		return err
	}

	f, err := atomicfile.New(string(filename), 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return encode(f, cfg)
}

func encode(w io.Writer, value any) error {
	buf, err := config.Marshal(value)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

func Load(filename ConfigPath) (*config.Config, error) {
	var cfg config.Config
	err := ReadConfigFile(filename, &cfg)
	if err != nil {
		return nil, err
	}
	// inject lotus api info into env
	lotusApiInfo := fmt.Sprintf("%s:%s", cfg.LotusInfo.Token, cfg.LotusInfo.APIAddress)
	if err := os.Setenv("FULLNODE_API_INFO", lotusApiInfo); err != nil {
		return nil, fmt.Errorf("failed to set FULLNODE_API_INFO into os environment variable: %s", err)
	}

	return &cfg, err
}
