package main

import (
	"fil-chain-extractor/config"
	"fil-chain-extractor/config/serialize"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"fil-chain-extractor/pkg/util/filesys"
)

func NewInitCmd() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initializes fce(fil-chain-extractor) config file.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "repo-path",
				Usage:    "repository directory persistent fce config file and data",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			repoPath := ctx.String("repo-path")
			var cfgPath string
			if repoPath == "" {
				cfgPath = filepath.Join(config.DefaultRepoRoot, config.DefaultConfigFile)
			} else {
				cfgPath = filepath.Join(repoPath, config.DefaultConfigFile)
			}
			cfgFileExist, err := filesys.IsFileExists(cfgPath)
			if err != nil {
				return fmt.Errorf("failure to check config file: %s", err)
			}
			if cfgFileExist {
				return fmt.Errorf("config file `%s` exists", repoPath)
			} else {
				cfg := config.Init()
				if err := serialize.WriteConfigFile(cfgPath, cfg); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
