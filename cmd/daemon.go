package main

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	lotusAPI "github.com/filecoin-project/lotus/api"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"

	"fil-chain-extractor/pkg/util/multiaddress"
)

func NewDaemonCmd() *cli.Command {
	return &cli.Command{
		Name:  "daemon",
		Usage: "Start a fil-chain-extractor daemon process",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "lotus-api-token",
				Usage:    "lotus api auth token which has read permission",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "lotus-api-address",
				Usage:    "lotus full-node api address in multi-address manner",
				Value:    "/ip4/127.0.0.1/tcp/1234",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "repo-path",
				Usage:    "daemon will read/persistent configs and data from the repository path",
				Required: false,
			},
		},
		Action: func(cCtx *cli.Context) error {

			headers := http.Header{}
			if authToken := cCtx.String("api-token"); authToken != "" {
				headers["Authorization"] = []string{"Bearer " + authToken}
			}
			apiUrl := cCtx.String("api-url")

			apiUrlMnet, err := multiaddress.ToNetAddress(apiUrl)
			if err != nil {
				logger.Fatalf("multiaddress parse failed, %s", err)
			}

			var api lotusAPI.FullNodeStruct
			closer, err := jsonrpc.NewMergeClient(
				context.Background(),
				"ws://"+apiUrlMnet+"/rpc/v0",
				"Filecoin",
				[]interface{}{&api.Internal, &api.CommonStruct.Internal},
				headers,
			)
			if err != nil {
				logger.Fatalf("connecting with lotus failed: %s", err)
			}
			defer closer()

			tipset, err := api.ChainHead(context.Background())
			notify, err := api.ChainNotify(context.TODO())
			if err != nil {
				return err
			}

			if err != nil {
				log.Fatalf("calling chain head: %s", err)
			}
			logger.Infof("Current chain head is : %s", tipset.String())
			return nil
		},
	}
}
