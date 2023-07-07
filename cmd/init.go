package main

import (
	"fmt"
	"github.com/pando-project/fil-chain-extractor/config"
	"github.com/pando-project/fil-chain-extractor/config/serialize"
	"github.com/urfave/cli/v2"

	"github.com/pando-project/fil-chain-extractor/pkg/util/filesys"
)

func NewInitCmd() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initializes fce(github.com/pando-project/fil-chain-extractor) config file.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "repo-path",
				Usage:    "repository directory persistent fce config file and data",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			repoPath := ctx.String("repo-path")
			cfgPath, err := config.Filename(repoPath, "")
			if err != nil {
				return fmt.Errorf("failed to get config path: %s", err)
			}
			cfgFileExist, err := filesys.IsFileExists(cfgPath)
			if err != nil {
				return fmt.Errorf("failure to check config file: %s", err)
			}
			if cfgFileExist {
				return fmt.Errorf("config file `%s` exists", cfgPath)
			} else {
				cfg := config.Init()
				if err := serialize.WriteConfigFile(serialize.ConfigPath(cfgPath), cfg); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
