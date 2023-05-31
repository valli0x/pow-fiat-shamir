package cmd

import (
	"encoding/json"
	"net/http"
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

			mux := http.NewServeMux()
			mux.Handle("/fiat-shamir", fiatShamirHandler())

			logger.Error("Server error:", http.ListenAndServe(config.Address, mux))

			return nil
		},
	}
	return cmd
}

func fiatShamirHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		respondOk(w, nil)
	})
}

func respondOk(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if body == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.Encode(body)
	}
}

func RespondError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	type ErrorResponse struct {
		Errors []string `json:"errors"`
	}
	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}
