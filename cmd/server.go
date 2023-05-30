package cmd

import (
	"os"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
)

type ServerFlags struct {
	LogsFormat bool
}

var (
	serverFlags = &ServerFlags{}
)

func init() {
	server := ClientCmd()

	server.PersistentFlags().BoolVar(&serverFlags.LogsFormat, "log_json", false, "logs format")

	RootCmd.AddCommand(server)
}

func ServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "server",
		Short:        "Fiat-Shamir server",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// create logger
			logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
				Output:     os.Stdout,
				Level:      hclog.DefaultLevel,
				JSONFormat: serverFlags.LogsFormat,
			})
			logger.Info("starting fiat-shamir server")

			// mux := http.NewServeMux()
			// mux.Handle("/fiat-shamir", )

			return nil
		},
	}
	return cmd
}
