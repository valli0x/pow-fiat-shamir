package cmd

import (
	"fmt"
	"net/http"
	"os"
	"pow-fiat-shamir/sdk"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

type ClientFlags struct {
	LogsFormat bool
}

var (
	clientFlags = &ClientFlags{}
)

func init() {
	client := ClientCmd()

	client.PersistentFlags().BoolVar(&clientFlags.LogsFormat, "log_json", false, "logs format")

	RootCmd.AddCommand(client)
}

func ClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "client",
		Short:        "Fiat-Shamir client",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// create logger
			logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
				Output:     os.Stdout,
				Level:      hclog.DefaultLevel,
				JSONFormat: clientFlags.LogsFormat,
			})
			logger.Info("starting fiat-shamir client")

			suite := edwards25519.NewBlakeSHA256Ed25519() // (client & server)

			// serverAddress := os.Getenv("FIAT_SHAMIR_SERVER")

			client := HttpClient()

			round1, err := StartRequest(client)
			if err != nil {
				return err
			}

			x := sdk.ComputeX(round1.Message, suite)
			xG, xH := sdk.ComputexGxH(suite, round1.G, round1.H, x)
			c := sdk.ComputeC(suite, round1.G, round1.H, xG, xH)
			v, vG, vH := sdk.ComputeVvGvH(suite, round1.G, round1.H)
			_, rG, rH := sdk.ComputeRrGrH(suite, c, x, v, round1.G, round1.H)

			round2 := &sdk.Round2{
				C: c, 

				XH: xH,
				XG: xG,

				RG: rG,
				RH: rH, 

				VG: vG,
				VH: vH,
			}

			quotes, err := ResultRequest(client, round2)
			if err != nil {
				return err
			}

			fmt.Println("Quotes from “word of wisdom” book:", quotes)

			return nil
		},
	}
	return cmd
}

func StartRequest(client *http.Client) (*sdk.Round1, error) {
	return nil, nil
}

func ResultRequest(client *http.Client, round2 *sdk.Round2) (string, error) {
	return "", nil
}

func HttpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}
