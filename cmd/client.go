package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
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

			// serverAddress := os.Getenv("FIAT_SHAMIR_SERVER")

			client := HttpClient()

			err = StartRequest(client)
			if err != nil {
				return err
			}

			quotes, err := ResultRequest(client)
			if err != nil {
				return err
			}

			fmt.Println("Quotes from “word of wisdom” book:", quotes)

			return nil
		},
	}
	return cmd
}

func StartRequest(client *http.Client) error {
	return nil
}

func ResultRequest(client *http.Client) (string, error) {
	return "", nil
}

func HttpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}
